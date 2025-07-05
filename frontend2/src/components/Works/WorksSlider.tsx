import React, { useEffect, useState, useRef } from "react";
import { useSelector } from "react-redux";
import styles from "./WorksSlider.module.css";
import { RootState } from "../../store/store";
import { getWorks } from "./Works.thunks";
import { useAppDispatch } from "../../hooks/hooks";
import WorkPreview from "./WorkPreview";
import { useTranslation } from "react-i18next";

type WorksSliderProps = {
     userId: string;
     interval?: number;
};

const WorksSlider: React.FC<WorksSliderProps> = ({ userId, interval = 5000 }) => {
     const dispatch = useAppDispatch();
     const { t } = useTranslation();

     const worksIds = useSelector((state: RootState) => state.works.worksIds);
     const loading = useSelector((state: RootState) => state.works.loading);
     const error = useSelector((state: RootState) => state.works.error);

     const [currentIndex, setCurrentIndex] = useState(0);
     const [previewImage, setPreviewImage] = useState<string | null>(null);
     const timeoutRef = useRef<NodeJS.Timeout | null>(null);

     useEffect(() => {
          dispatch(getWorks(userId));
     }, [dispatch, userId]);

     const resetTimeout = () => {
          if (timeoutRef.current) {
               clearTimeout(timeoutRef.current);
          }
     };

     useEffect(() => {
          if (!worksIds || worksIds.length === 0) return;

          resetTimeout();
          timeoutRef.current = setTimeout(() => {
               setCurrentIndex((prevIndex) =>
                    prevIndex === worksIds.length - 1 ? 0 : prevIndex + 1
               );
          }, interval);

          return () => {
               resetTimeout();
          };
     }, [currentIndex, worksIds, interval]);

     const prevSlide = () => {
          resetTimeout();
          setCurrentIndex((prev) =>
               prev === 0 ? worksIds.length - 1 : prev - 1
          );
     };

     const nextSlide = () => {
          resetTimeout();
          setCurrentIndex((prev) =>
               prev === worksIds.length - 1 ? 0 : prev + 1
          );
     };

     const openPreview = (workId: string) => {
          setPreviewImage(`${__BASE_API_URL__}/users/${userId}/works/${workId}`);
     };

     const closePreview = () => {
          setPreviewImage(null);
     };

     if (loading) return <div>{t("worksSlider.loading")}</div>;
     if (error) return <div>{t("worksSlider.error", { error })}</div>;
     if (!worksIds || worksIds.length === 0) return null;

     return (
          <div className={styles.worksSection}>
               <div className={styles.slider}>
                    <button
                         onClick={prevSlide}
                         className={styles.navButton}
                         aria-label={t("worksSlider.prev")}
                    >
                         ‹
                    </button>
                    <div className={styles.workItem}>
                         {worksIds.map((work, index) => (
                              <img
                                   key={work}
                                   src={`${__BASE_API_URL__}/users/${userId}/works/${work}`}
                                   alt={t("worksSlider.workAlt", { index: index + 1 })}
                                   className={`${styles.workImage} ${index === currentIndex ? styles.active : styles.inactive
                                        }`}
                                   loading="lazy"
                                   onClick={() => openPreview(work.toString())}
                              />
                         ))}
                    </div>
                    <button
                         onClick={nextSlide}
                         className={styles.navButton}
                         aria-label={t("worksSlider.next")}
                    >
                         ›
                    </button>
               </div>
               {previewImage && (
                    <WorkPreview imageUrl={previewImage} onClose={closePreview} />
               )}
          </div>
     );
};

export default WorksSlider;
