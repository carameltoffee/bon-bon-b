import { ThunkAction } from "redux-thunk";
import { RootState } from "../../store/store";
import { AuthActionTypes, loginFailure, request, loginSuccess, registerFailure, registerSuccess, logoutAction, userLoaded } from "./Auth.actions";
import { setErrorAlert, setSuccessAlert } from "../Alert/Alert.thunks";
import { isJwtExpired } from "../../utils/jwt";

export type AppThunk<ReturnType = void> = ThunkAction<
     ReturnType,
     RootState,
     unknown,
     AuthActionTypes
>;

export const loadUser = (): AppThunk => async (dispatch) => {
     const userJSON = localStorage.getItem('user') || '{}';
     const userTokenJSON = localStorage.getItem('token') || '{}';
     const token = JSON.parse(userTokenJSON);
     const user = JSON.parse(userJSON) as IUser;
     if (isJwtExpired(token)) {
          dispatch(logout());
          return;
     }
     dispatch(userLoaded(token, user))
}

export const loginUser = (credentials: ILoginReq): AppThunk => async (dispatch) => {
     dispatch(request());

     try {
          const response = await fetch(`${__BASE_API_URL__}/login`, {
               method: 'POST',
               headers: { 'Content-Type': 'application/json' },
               body: JSON.stringify(credentials),
          });

          const data = await response.json();
          if (!response.ok) throw new Error(data?.error || 'Неизвестная ошибка');

          const { payload } = decodeJWT(data.token);
          const responseGetUser = await fetch(`${__BASE_API_URL__}/users/${payload.id}`, {
               method: 'GET',
               headers: { 'Content-Type': 'application/json' },
          });
          const userData = await responseGetUser.json();
          if (!responseGetUser.ok) throw new Error(data?.error || 'Неизвестная ошибка');
          dispatch(loginSuccess(userData, data.token));
          dispatch(setSuccessAlert("Вы успешно залогинились!"));
     } catch (error) {
          if (error instanceof Error) {
               dispatch(setErrorAlert(error.message));
               dispatch(loginFailure(error.message));
          } else {
               dispatch(loginFailure('Неизвестная ошибка'));
          }
     }
};

export const registerUser = (req: IRegisterReq): AppThunk => async (dispatch) => {
     dispatch(request());

     try {
          const response = await fetch(`${__BASE_API_URL__}/register`, {
               method: 'POST',
               headers: { 'Content-Type': 'application/json' },
               body: JSON.stringify(req),
          });

          const data = await response.json();
          if (!response.ok) throw new Error(data?.error || 'Неизвестная ошибка');

          dispatch(registerSuccess());
          dispatch(setSuccessAlert("Вы успешно зарегестрировались! Теперь войдите."))
     } catch (error) {
          if (error instanceof Error) {
               dispatch(setErrorAlert(error.message));
               dispatch(registerFailure(error.message));
          } else {
               dispatch(registerFailure('Неизвестная ошибка'));
          }
     }
};

export const sendCode = (email: string): AppThunk => async (dispatch) => {
     try {
          const response = await fetch(`${__BASE_API_URL__}/send-code`, {
               method: 'POST',
               headers: { 'Content-Type': 'application/json' },
               body: JSON.stringify({ email }),
          });

          if (!response.ok) throw new Error('Ошибка отправки кода');
          dispatch(setSuccessAlert('Код отправлен на почту'));
     } catch (error) {
          const message = error instanceof Error ? error.message : 'Неизвестная ошибка';
          dispatch(setErrorAlert(message));
     }
};

export const restore = (email: string, code: string, password: string): AppThunk => async (dispatch) => {
     try {
          const response = await fetch(`${__BASE_API_URL__}/restore`, {
               method: 'POST',
               headers: { 'Content-Type': 'application/json' },
               body: JSON.stringify({ email, code, password }),
          });
          if (!response.ok) {
               const data = await response.json();
               throw new Error(data?.error || 'Невозможно изменить пароль')
          };
          dispatch(registerSuccess());
          dispatch(setSuccessAlert('Успешно изменен пароль для входа'));
     } catch (error) {
          const message = error instanceof Error ? error.message : 'Неизвестная ошибка';
          dispatch(setErrorAlert(message));
          dispatch(registerFailure(message));
     }
}

export const logout = (): AppThunk => async (dispatch) => {
     dispatch(logoutAction());
};

function decodeJWT(token: string): { payload?: any } {
     if (!token) return {};
     const parts = token.split('.');
     if (parts.length !== 3) return {};

     try {
          const payloadJson = atob(parts[1]);
          const payload = JSON.parse(payloadJson);
          return { payload };
     } catch {
          return {};
     }
}
