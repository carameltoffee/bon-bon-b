import asyncio
import json
import aio_pika
import logging
import os
import smtplib
from email.mime.text import MIMEText

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

smtp_server = None
rmq_connection = None
rmq_channel = None

async def main():
    global smtp_server, rmq_connection, rmq_channel
    try:
        smtp_server = smtplib.SMTP(os.getenv('SMTP_HOST'), os.getenv('SMTP_PORT'))
        smtp_server.starttls()
        smtp_server.login(os.getenv('SMTP_LOGIN'), os.getenv('SMTP_PASSWORD'))
        
        rmq_url = os.getenv('RABBITMQ_URL')
        rmq_connection = await aio_pika.connect_robust(rmq_url)
        rmq_channel = await rmq_connection.channel()
        await rmq_channel.set_qos(prefetch_count=10)

        exchanges = {
            "user": aio_pika.ExchangeType.TOPIC,
            "book": aio_pika.ExchangeType.TOPIC
        }

        queues = {}
        bindings = {
            "user": ["user.logged"],
            "book": ["book.booked"]
        }

        for ex_name, ex_type in exchanges.items():
            exchange = await rmq_channel.declare_exchange(ex_name, ex_type, durable=True)
            queue = await rmq_channel.declare_queue(f"{ex_name}_queue", durable=True)
            queues[ex_name] = queue

            for rk in bindings[ex_name]:
                await queue.bind(exchange, rk)

        logger.info("Starting consumers...")
        await asyncio.gather(
               queues["user"].consume(lambda msg: on_message(msg, "user")),
               queues["book"].consume(lambda msg: on_message(msg, "book")),
               asyncio.Future()
          )

    except Exception as e:
        logger.error(f"Main function error: {e}")
    finally:
        if rmq_connection:
          await rmq_connection.close()
        if smtp_server:
            smtp_server.quit()
            
def build_email(to, subject, body):
    msg = MIMEText(body)
    msg['Subject'] = subject
    msg['To'] = to
    msg['From'] = os.getenv('SMTP_FROM', 'noreply@example.com')
    return msg

async def on_message(message: aio_pika.IncomingMessage, exchange_name):
    async with message.process():
        try:
            payload = json.loads(message.body)
            email = payload.get('email', None)
            if not email:
                logger.error("Email not found in message payload")
                return

            routing_key = message.routing_key
            msg = None

            match exchange_name:
                case "user":
                    if routing_key == "user.logged":
                        msg = build_email(email, f"Вход в аккаунт / {os.getenv("APP_NAME")}", "Кто-то пытался войти в ваш аккаунт")
                    else:
                        logger.warning(f"Unhandled routing key for user exchange: {routing_key}")
                        return
                case "book":
                    if routing_key == "book.booked":
                        slot = payload.get('slot', 'Не указано')
                        contact_phone = payload.get('contact_phone', 'Не указано')
                        msg = build_email(email, f"Запись! / {os.getenv('APP_NAME', 'App')}", f"К вам сделали бронь на {slot}, контактный номер: {contact_phone}")
                    else:
                        logger.warning(f"Unhandled routing key for book exchange: {routing_key}")
                        return
                case _:
                    logger.warning(f"Unhandled exchange name: {exchange_name}")
                    return

            if msg:
                smtp_server.send_message(msg)
        except Exception as e:
            logger.exception(f"Cannot process the message: {e}")

if __name__ == '__main__':
    asyncio.run(main())