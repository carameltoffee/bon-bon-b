import { Alert } from '../components/Alert/Alert';
import { InitRoutes } from '../routes/Routes';
import './App.css';
import { useEffect } from 'react';
import { loadUser } from '../components/Auth/Auth.thunks';
import Header from '../components/Header/Header';
import { BrowserRouter } from 'react-router-dom';
import { useAppDispatch } from '../hooks/hooks';
import { ConfirmProvider } from '../hooks/Confirm/Confirm';
import { useAutoLanguage } from '../hooks/useAutoLang';

function App() {
	const dispatch = useAppDispatch();

	useAutoLanguage();

	useEffect(() => {
		dispatch(loadUser());
	}, [dispatch]);

	return (
		<>
			<BrowserRouter>
				<ConfirmProvider>
					<Header />
					<InitRoutes />
					<Alert />
				</ConfirmProvider>
			</BrowserRouter>
		</>
	);
}

export default App;
