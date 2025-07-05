import React, { useEffect } from "react";
import { useAppDispatch } from "../../hooks/hooks";
import { useSelector } from "react-redux";
import { RootState } from "../../store/store";
import { GetAppointments, DeleteAppointment } from "./Appointments.thunks";
import { GetMasterById } from "../Masters/Masters.thunks";
import AppointmentCard from "./AppointmentCard";
import { useConfirm } from "../../hooks/Confirm/Confirm";
import { Spinner } from "../Spinner/Spinner";
import { useTranslation } from "react-i18next";

const UserAppointmentsList: React.FC = () => {
     const dispatch = useAppDispatch();
     const { items, loading, error } = useSelector((state: RootState) => state.appointments);
     const mastersById = useSelector((state: RootState) => state.masters.mastersById);
     const token = useSelector((state: RootState) => state.auth.token);
     const confirm = useConfirm();
     const { t } = useTranslation();

     useEffect(() => {
          if (!token) return;
          dispatch(GetAppointments(token));
     }, [dispatch, token]);

     useEffect(() => {
          if (!token) return;
          items.user_appointments.forEach(appt => {
               if (!mastersById[appt.master_id]) {
                    dispatch(GetMasterById(appt.master_id.toString()));
               }
          });
          items.master_appointments.forEach(appt => {
               if (!mastersById[appt.user_id]) {
                    dispatch(GetMasterById(appt.user_id.toString()));
               }
          });
     }, [items, mastersById, dispatch, token]);

     const handleDelete = (id: number) => {
          if (!token) return;
          confirm({
               message: t("appointments.delete_confirm"),
               onConfirm: () => dispatch(DeleteAppointment(token, id)),
          });
     };

     if (loading) return <Spinner />;
     if (error) return <p className="text-red-500">{error}</p>;
     if (!items) return null;

     if (!items.user_appointments.length && !items.master_appointments.length)
          return <p>{t("appointments.no_appointments")}</p>;

     return (
          <div className="space-y-8">
               {!!items.user_appointments.length && (
                    <div>
                         <h2 className="text-xl font-semibold mb-2">{t("appointments.title_user")}</h2>
                         <ul className="space-y-3">
                              {items.user_appointments.map((appt) => {
                                   const master = mastersById[appt.master_id];
                                   return (
                                        <AppointmentCard
                                             key={appt.id}
                                             id={appt.id}
                                             spec={master ? master.specialization : t("appointments.loading_spec")}
                                             name={
                                                  master
                                                       ? master.full_name
                                                       : t("appointments.loading_name", { id: appt.master_id })
                                             }
                                             date={appt.scheduled_at}
                                             onDelete={handleDelete}
                                             variant="user"
                                        />
                                   );
                              })}
                         </ul>
                    </div>
               )}

               {!!items.master_appointments.length && (
                    <div>
                         <h2 className="text-xl font-semibold mb-2">{t("appointments.title_master")}</h2>
                         <ul className="space-y-3">
                              {items.master_appointments.map((appt) => {
                                   const user = mastersById[appt.user_id];
                                   return (
                                        <AppointmentCard
                                             key={appt.id}
                                             id={appt.id}
                                             name={
                                                  user
                                                       ? user.full_name
                                                       : t("appointments.loading_name", { id: appt.user_id })
                                             }
                                             date={appt.scheduled_at}
                                             onDelete={handleDelete}
                                             variant="master"
                                        />
                                   );
                              })}
                         </ul>
                    </div>
               )}
          </div>
     );
};

export default UserAppointmentsList;
