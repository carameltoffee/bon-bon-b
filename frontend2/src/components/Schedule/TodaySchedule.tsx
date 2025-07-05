import React, { useEffect } from "react";
import { useSelector } from "react-redux";
import { RootState } from "../../store/store";
import { GetSchedule } from "./Schedule.thunks";
import { useAppDispatch } from "../../hooks/hooks";
import { Spinner } from "../Spinner/Spinner";
import styles from "./TodaySchedule.module.css";
import { useTranslation } from "react-i18next";

interface TodayScheduleProps {
     userId: string;
}

const TodaySchedule: React.FC<TodayScheduleProps> = ({ userId }) => {
     const dispatch = useAppDispatch();
     const { t } = useTranslation();
     const { loading, error, schedule } = useSelector((state: RootState) => state.schedule);
     const today = new Date().toISOString().split("T")[0];

     useEffect(() => {
          dispatch(GetSchedule(today, userId));
     }, [dispatch, today, userId]);

     if (loading) return <Spinner />;
     if (error) return <div>{t("todaySchedule.error", { error })}</div>;
     if (!schedule) return null;
     if (schedule.days_off.includes(today))
          return (
               <div className={styles.container}>
                    <h2 className={styles.header}>{t("todaySchedule.dayOff")}</h2>
               </div>
          );

     const appointments = schedule.appointments;

     if (!appointments || appointments.length === 0) {
          return (
               <div className={styles.container}>
                    <h2 className={styles.header}>{t("todaySchedule.noAppointments")}</h2>
               </div>
          );
     }

     return (
          <div className={styles.container}>
               <h2 className={styles.header}>
                    {t("todaySchedule.appointmentsCount", { count: appointments.length })}
               </h2>
               <div className={styles.cards}>
                    {appointments.map((time, i) => (
                         <div key={i} className={styles.card}>
                              <div className={styles.time}>{time}</div>
                         </div>
                    ))}
               </div>
          </div>
     );
};

export default TodaySchedule;
