import React, { useState } from 'react';
import { useAppDispatch } from '../../hooks/hooks';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/store';
import { setWorkingHoursByWeekday } from './Schedule.thunks';
import styles from './ScheduleEditor.module.css';
import { useTranslation } from 'react-i18next';

const weekdays = [
     'monday',
     'tuesday',
     'wednesday',
     'thursday',
     'friday',
     'saturday',
     'sunday',
];

const ScheduleByWeekday: React.FC = () => {
     const dispatch = useAppDispatch();
     const { t } = useTranslation();
     const token = useSelector((state: RootState) => state.auth.token);
     if (!token) return null;

     const [weekday, setWeekday] = useState(0);
     const [weekdaySlots, setWeekdaySlots] = useState<string[]>([]);
     const [weekdaySlotInput, setWeekdaySlotInput] = useState('');

     const addSlot = (
          input: string,
          setSlots: React.Dispatch<React.SetStateAction<string[]>>,
          slots: string[],
          setInput: React.Dispatch<React.SetStateAction<string>>
     ) => {
          if (input && !slots.includes(input)) {
               setSlots([...slots, input].sort());
               setInput('');
          }
     };

     const removeSlot = (
          slot: string,
          slots: string[],
          setSlots: React.Dispatch<React.SetStateAction<string[]>>
     ) => {
          setSlots(slots.filter(s => s !== slot));
     };

     const handleSaveWeekday = () => {
          const dayEn = weekdays[weekday];
          if (weekdaySlots.length > 0) {
               dispatch(setWorkingHoursByWeekday(token, dayEn, weekdaySlots));
          }
     };

     return (
          <div className={styles.tabContent}>
               <label className={styles.label}>{t('scheduleByWeekday.selectWeekday')}</label>
               <select
                    className={styles.input}
                    value={weekday}
                    onChange={(e) => setWeekday(Number(e.target.value))}
               >
                    {weekdays.map((day, idx) => (
                         <option key={day} value={idx}>
                              {t(`weekdays.${day}`)}
                         </option>
                    ))}
               </select>

               <div className={styles.slotRow}>
                    <input
                         type="time"
                         value={weekdaySlotInput}
                         onChange={(e) => setWeekdaySlotInput(e.target.value)}
                         className={styles.input}
                    />
                    <button
                         className={styles.buttonPrimary}
                         onClick={() =>
                              addSlot(weekdaySlotInput, setWeekdaySlots, weekdaySlots, setWeekdaySlotInput)
                         }
                    >
                         {t('scheduleByWeekday.addSlot')}
                    </button>
               </div>

               <div className={styles.slots}>
                    {weekdaySlots.map(slot => (
                         <span
                              key={slot}
                              className={styles.slot}
                              onClick={() => removeSlot(slot, weekdaySlots, setWeekdaySlots)}
                         >
                              {slot} ×
                         </span>
                    ))}
               </div>

               <div className={styles.section}>
                    <button className={styles.buttonSuccess} onClick={handleSaveWeekday}>
                         {t('scheduleByWeekday.saveFor', { day: t(`weekdays.${weekdays[weekday]}`) })}
                    </button>
               </div>
          </div>
     );
};

export default ScheduleByWeekday;
