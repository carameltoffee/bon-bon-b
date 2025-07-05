import React, { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { RootState } from "../../store/store";
import { useAppDispatch } from "../../hooks/hooks";
import { getWorks, deleteWork, addWork } from "./Works.thunks";
import styles from "./WorksEditor.module.css";
import { useTranslation } from "react-i18next";

const WorksEditor: React.FC = () => {
     const dispatch = useAppDispatch();
     const { worksIds, loading, error } = useSelector((state: RootState) => state.works);
     const token = useSelector((state: RootState) => state.auth.token);
     const user = useSelector((state: RootState) => state.auth.user);
     const { t } = useTranslation();

     const [file, setFile] = useState<File | null>(null);

     useEffect(() => {
          if (user) {
               dispatch(getWorks(user.id.toString()));
          }
     }, [dispatch, user]);

     const handleUpload = () => {
          if (file && token && user) {
               dispatch(addWork(token, file));
               setFile(null);
               dispatch(getWorks(user.id.toString()));
          }
     };

     const handleDelete = (id: number) => {
          if (token && user) {
               dispatch(deleteWork(token, id));
               dispatch(getWorks(user.id.toString()));
          }
     };

     if (!token || !user) {
          return <p>{t("worksEditor.noAccess")}</p>;
     }

     return (
          <div className={styles.container}>
               <div className={styles.uploadBlock}>
                    <input
                         type="file"
                         accept="image/*"
                         onChange={(e) => setFile(e.target.files?.[0] || null)}
                         className={styles.fileInput}
                    />
                    <button
                         onClick={handleUpload}
                         disabled={!file}
                         className={styles.uploadButton}
                    >
                         {t("worksEditor.upload")}
                    </button>
               </div>

               {error && <p className="text-red-500">{error}</p>}

               <div className={styles.grid}>
                    {loading && <p>{t("worksEditor.loading")}</p>}
                    {!loading && (!worksIds || worksIds.length === 0) && (
                         <p className={styles.message}>{t("worksEditor.noWorks")}</p>
                    )}
                    {!loading &&
                         worksIds &&
                         worksIds.map((id) => (
                              <div key={id} className={styles.card}>
                                   <img
                                        src={`${__BASE_API_URL__}/users/${user.id}/works/${id}`}
                                        alt={t("worksEditor.workAlt", { id })}
                                        className={styles.image}
                                   />
                                   <button
                                        onClick={() => handleDelete(id)}
                                        className={styles.deleteButton}
                                        aria-label={t("worksEditor.workAlt", { id })}
                                   >
                                        ✕
                                   </button>
                              </div>
                         ))}
               </div>
          </div>
     );
};

export default WorksEditor;
