import React, { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { RootState } from "../../store/store";
import { Spinner } from "../Spinner/Spinner";
import { useAppDispatch } from "../../hooks/hooks";
import { useNavigate } from "react-router-dom";
import { searchUsers } from "../Search/Search.thunks";
import debounce from "lodash.debounce";
import styles from "./MastersList.module.css";
import Stars from "../Stars/Stars";
import { useTranslation } from "react-i18next";

const MastersList: React.FC = () => {
     const { t } = useTranslation();
     const dispatch = useAppDispatch();
     const navigate = useNavigate();

     const { users, loading, error } = useSelector((state: RootState) => state.search);
     const safeUsers = users || [];
     const [search, setSearch] = useState("");

     const debouncedSearch = React.useMemo(
          () =>
               debounce((q: string) => {
                    if (q.trim()) {
                         dispatch(searchUsers(q.trim()));
                    }
               }, 400),
          [dispatch]
     );

     useEffect(() => {
          dispatch(searchUsers(""));
     }, [dispatch]);

     useEffect(() => {
          if (search.trim()) {
               debouncedSearch(search);
          } else {
               dispatch(searchUsers(""));
          }

          return () => {
               debouncedSearch.cancel();
          };
     }, [search, debouncedSearch, dispatch]);

     return (
          <div className={styles.wrapper}>
               <input
                    type="text"
                    placeholder={t("mastersList.searchPlaceholder")}
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    className={styles.input}
               />

               {loading && <Spinner />}
               {!loading && error && (
                    <p style={{ color: "red" }}>
                         {t("mastersList.searchError", { error })}
                    </p>
               )}
               {!loading && !error && safeUsers.length === 0 && (
                    <p style={{ color: "#777" }}>
                         {t("mastersList.noMastersFound")}
                    </p>
               )}

               <div className={styles.grid}>
                    {safeUsers.map((master: IUser) => (
                         <div
                              key={master.id}
                              onClick={() => navigate(`/res/${master.id}`)}
                              className={styles.card}
                         >
                              <h2 className={styles.cardTitle}>{master.full_name}</h2>
                              <p className={styles.cardText}>{master.specialization}</p>
                              {master.average_rating !== undefined && (
                                   <Stars rating={master.average_rating} />
                              )}
                         </div>
                    ))}
               </div>
          </div>
     );
};

export default MastersList;
