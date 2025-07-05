import React, { useState, useEffect } from 'react';
import Calendar from 'react-calendar';
import { useAppDispatch } from '../../hooks/hooks';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/store';
import {
     setDayOff,
     deleteDayOff,
     setWorkingHoursByDate,
     deleteWorkingHoursByDate,
     GetSchedule,
} from './Schedule.thunks';
import styles from './ScheduleEditor.module.css';
import { useTranslation } from 'react-i18next';

const ScheduleByDate: React.FC = () => {
     const dispatch = useAppDispatch();
     const { t } = useTranslation();

     const token = useSelector((state: RootState) => state.auth.token);
     const userId = useSelector((state: RootState) => state.auth.user?.id);
     const schedule = useSelector((state: RootState) => state.schedule.schedule);

     const [selectedDate, setSelectedDate] = useState<Date | null>(null);
     const [dateSlots, setDateSlots] = useState<string[]>([]);
     const [dateSlotInput, setDateSlotInput] = useState('');
     const [isDayOff, setIsDayOff] = useState(false);

     if (!token || !userId) return null;

     useEffect(() => {
          if (!selectedDate) return;

          const isoDate = selectedDate.toLocaleDateString('sv-SE');
          dispatch(GetSchedule(isoDate, userId.toString()));
     }, [selectedDate, userId, dispatch]);

     useEffect(() => {
          if (!selectedDate || !schedule) return;

          const isoDate = selectedDate.toLocaleDateString('sv-SE');

          const dayIsOff = schedule.days_off?.includes(isoDate) ?? false;
          setIsDayOff(dayIsOff);

          if (dayIsOff) {
               setDateSlots([]);
          } else {
               setDateSlots(schedule.slots ?? []);
          }
     }, [schedule, selectedDate]);

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

     const handleSaveDate = () => {
          if (!selectedDate) return;
          const isoDate = selectedDate.toLocaleDateString('sv-SE');

          if (isDayOff) {
               dispatch(setDayOff(token, isoDate));
          } else {
               dispatch(deleteDayOff(token, isoDate));
               if (dateSlots.length === 0) {
                    dispatch(deleteWorkingHoursByDate(token, isoDate));
               } else {
                    dispatch(setWorkingHoursByDate(token, isoDate, dateSlots));
               }
          }
     };

     const handleRemoveDayOff = () => {
          if (!selectedDate) return;
          dispatch(deleteDayOff(token, selectedDate.toLocaleDateString('sv-SE')));
     };

     return (
          <div className={styles.tabContent}>
               <Calendar
                    onChange={(value) => {
                         if (value instanceof Date) setSelectedDate(value);
                    }}
                    value={selectedDate}
                    selectRange={false}
                    className={styles.calendar}
               />

               {selectedDate && (
                    <div>
                         <label className={styles.switch}>
                              <input
                                   type="checkbox"
                                   checked={isDayOff}
                                   onChange={(e) => setIsDayOff(e.target.checked)}
                              />
                              <span className={styles.slider}></span>
                              <span>{t('scheduleByDate.dayOff')}</span>
                         </label>

                         {!isDayOff && (
                              <>
                                   <div className={styles.slotRow}>
                                        <input
                                             type="time"
                                             value={dateSlotInput}
                                             onChange={(e) => setDateSlotInput(e.target.value)}
                                             className={styles.input}
                                        />
                                        <button
                                             className={styles.buttonPrimary}
                                             onClick={() => addSlot(dateSlotInput, setDateSlots, dateSlots, setDateSlotInput)}
                                        >
                                             {t('scheduleByDate.addSlot')}
                                        </button>
                                   </div>

                                   <div className={styles.slots}>
                                        {dateSlots.map(slot => (
                                             <span
                                                  key={slot}
                                                  className={styles.slot}
                                                  onClick={() => removeSlot(slot, dateSlots, setDateSlots)}
                                             >
                                                  {slot} ×
                                             </span>
                                        ))}
                                   </div>
                              </>
                         )}

                         <div className={styles.section}>
                              <button className={styles.buttonSuccess} onClick={handleSaveDate}>
                                   {t('scheduleByDate.save')}
                              </button>

                              {isDayOff && (
                                   <button className={styles.buttonDanger} onClick={handleRemoveDayOff}>
                                        {t('scheduleByDate.deleteDayOff')}
                                   </button>
                              )}
                         </div>
                    </div>
               )}
          </div>
     );
};

export default ScheduleByDate;
