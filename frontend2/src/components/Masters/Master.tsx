import React, { useEffect } from "react";
import { useSelector } from "react-redux";
import { RootState } from "../../store/store";
import { GetMasterById } from "./Masters.thunks";
import Avatar from "../Avatar/Avatar";
import { useAppDispatch } from "../../hooks/hooks";
import WorksSlider from "../Works/WorksSlider";
import styles from "./Master.module.css";
import Stars from "../Stars/Stars";
import { useTranslation } from "react-i18next";

type MasterProps = {
	masterId: string;
};

const Master: React.FC<MasterProps> = ({ masterId }) => {
	const { t } = useTranslation();
	const dispatch = useAppDispatch();

	const master = useSelector(
		(state: RootState) => state.masters.mastersById[masterId]
	);
	const isLoading = useSelector((state: RootState) => state.masters.loading);
	const error = useSelector((state: RootState) => state.masters.error);

	useEffect(() => {
		if (!master) {
			dispatch(GetMasterById(masterId));
		}
	}, [dispatch, masterId, master]);

	if (isLoading) return <div>{t("master.loading")}</div>;
	if (error || !master)
		return <div style={{ color: "red" }}>{t("master.error", { error })}</div>;

	return (
		<div className={styles.container}>
			<div className={styles.leftCol}>
				<div className={styles.avatar}>
					<Avatar userId={masterId} />
				</div>
				<h1>{master.full_name}</h1>
				<p>
					<strong>{t("master.username")}:</strong> {master.username}
				</p>
				{master.specialization && (
					<p>
						<strong>{t("master.specialization")}:</strong> {master.specialization}
					</p>
				)}
				<p>
					<strong>{t("master.registered")}:</strong>{" "}
					{new Date(master.registered_at).toLocaleDateString()}
				</p>
				{master.average_rating !== undefined && (
					<Stars rating={master.average_rating} />
				)}
			</div>

			<div className={styles.rightCol}>
				<WorksSlider userId={masterId} />
			</div>
		</div>
	);
};

export default Master;
