import React from "react";
import { formatTimestampUTC } from "../../utils/dates";
import styles from "./AppointmentCard.module.css";
import { useTranslation } from "react-i18next";

interface AppointmentCardProps {
     id: number;
     name: string;
     date: string;
     status?: string;
     spec?: string;
     onDelete: (id: number) => void;
     variant?: "user" | "master";
}

const AppointmentCard: React.FC<AppointmentCardProps> = ({
     id,
     name,
     spec,
     date,
     status,
     onDelete,
     variant = "user",
}) => {
     const { t } = useTranslation();
     const bgClass = variant === "user" ? styles.bgUser : styles.bgMaster;

     return (
          <li className={`${styles.card} ${bgClass}`}>
               <div className={styles.data}>
                    <div>
                         <strong>
                              {variant === "user"
                                   ? t("appointments.master")
                                   : t("appointments.client")}
                              :
                         </strong>{" "}
                         {name}
                         {spec && ` (${spec})`}
                    </div>
                    <div>
                         <strong>{t("appointments.date")}:</strong> {formatTimestampUTC(date)}
                    </div>
                    {status && (
                         <div>
                              <strong>{t("appointments.status")}:</strong> {status}
                         </div>
                    )}
               </div>
               <button className={styles.deleteBtn} onClick={() => onDelete(id)}>
                    {t("appointments.cancel")}
               </button>
          </li>
     );
};

export default AppointmentCard;
