import React, { useEffect } from "react";
import { useSelector } from "react-redux";
import { RootState } from "../../store/store";
import { useAppDispatch } from "../../hooks/hooks";
import { GetReviews } from "./Reviews.thunks";
import Avatar from "../Avatar/Avatar";
import { useTranslation } from "react-i18next";
import styles from "./ReviewsList.module.css";

type ReviewsListProps = {
     masterId: string;
};

const ReviewsList: React.FC<ReviewsListProps> = ({ masterId }) => {
     const dispatch = useAppDispatch();
     const { t } = useTranslation();

     const { reviewsByMaster, loading, error } = useSelector((state: RootState) => state.reviews);

     useEffect(() => {
          dispatch(GetReviews(masterId));
     }, [dispatch, masterId]);

     if (loading) return <p>{t("reviewsList.loading")}</p>;
     if (error) return null;

     return (
          <div className={styles.container}>
               <h2 className={styles.title}>{t("reviewsList.title")}</h2>
               {!reviewsByMaster[masterId] ? (
                    <p className={styles.noReviews}>{t("reviewsList.noReviews")}</p>
               ) : (
                    <ul className={styles.list}>
                         {reviewsByMaster[masterId].map((review) => (
                              <li key={review.id} className={styles.listItem}>
                                   <Avatar userId={review.user_id.toString()} size={30} />
                                   <div className={styles.reviewContent}>
                                        <div className={styles.rating}>
                                             {t("reviewsList.rating")}: {review.rating}/5
                                        </div>
                                        <p className={styles.comment}>{review.comment}</p>
                                   </div>
                              </li>
                         ))}
                    </ul>
               )}
          </div>
     );
};

export default ReviewsList;
