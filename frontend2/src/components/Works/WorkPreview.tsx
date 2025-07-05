import React from "react";
import { useTranslation } from "react-i18next";
import styles from "./WorkPreview.module.css";

type WorkPreviewProps = {
     imageUrl: string;
     onClose: () => void;
};

const WorkPreview: React.FC<WorkPreviewProps> = ({ imageUrl, onClose }) => {
     const { t } = useTranslation();

     return (
          <div className={styles.overlay} onClick={onClose}>
               <div className={styles.previewContent}>
                    <img
                         src={imageUrl}
                         alt={t("workPreview.altText")}
                         className={styles.previewImage}
                    />
                    <button
                         className={styles.closeButton}
                         onClick={onClose}
                         aria-label={t("workPreview.closeLabel")}
                    >
                         ×
                    </button>
               </div>
          </div>
     );
};

export default WorkPreview;
