import React, { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { RootState } from "../../store/store";
import { GetMasterById, UpdateMaster } from "./Masters.thunks";
import Avatar from "../Avatar/Avatar";
import { useAppDispatch } from "../../hooks/hooks";
import styles from "./Master.module.css";
import { setErrorAlert } from "../Alert/Alert.thunks";
import { useTranslation } from "react-i18next";

type MasterProps = {
     masterId: string;
};

const MasterEditable: React.FC<MasterProps> = ({ masterId }) => {
     const { t } = useTranslation();
     const dispatch = useAppDispatch();

     const master = useSelector(
          (state: RootState) => state.masters.mastersById[masterId]
     );
     const isLoading = useSelector((state: RootState) => state.masters.loading);
     const error = useSelector((state: RootState) => state.masters.error);
     const token = useSelector((state: RootState) => state.auth.token);

     const [form, setForm] = useState({
          id: 0,
          full_name: "",
          bio: "",
          username: "",
          email: "",
          specialization: "",
          registered_at: "",
     });

     useEffect(() => {
          if (!master) {
               dispatch(GetMasterById(masterId));
          } else {
               setForm({
                    id: Number(masterId),
                    full_name: master.full_name || "",
                    bio: master.bio || "",
                    username: master.username || "",
                    email: master.email || "",
                    specialization: master.specialization || "",
                    registered_at: master.registered_at,
               });
          }
     }, [dispatch, masterId, master]);

     const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
          setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
     };

     const handleSubmit = (e: React.FormEvent) => {
          e.preventDefault();
          if (!token) {
               dispatch(setErrorAlert(t("masterEditable.notAuthorized")));
               return;
          }
          dispatch(UpdateMaster(form, token));
     };

     if (isLoading) return <div>{t("masterEditable.loading")}</div>;
     if (error || !master)
          return <div style={{ color: "red" }}>{t("masterEditable.error", { error })}</div>;

     return (
          <div className={styles.container}>
               <Avatar userId={masterId} editable={true} />
               <form onSubmit={handleSubmit} className={styles.form}>
                    <label>
                         {t("masterEditable.fullName")}:
                         <input name="full_name" value={form.full_name} onChange={handleChange} />
                    </label>
                    <label>
                         {t("masterEditable.username")}:
                         <input name="username" value={form.username} onChange={handleChange} />
                    </label>
                    <label>
                         {t("masterEditable.email")}:
                         <input type="email" name="email" value={form.email} onChange={handleChange} />
                    </label>
                    <label>
                         {t("masterEditable.specialization")}:
                         <input name="specialization" value={form.specialization} onChange={handleChange} />
                    </label>
                    <button type="submit">{t("masterEditable.save")}</button>
               </form>
               <p>
                    <strong>{t("masterEditable.registered")}:</strong>{" "}
                    {new Date(master.registered_at).toLocaleDateString()}
               </p>
          </div>
     );
};

export default MasterEditable;
