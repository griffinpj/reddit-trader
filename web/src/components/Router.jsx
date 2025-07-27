import { BrowserRouter, Routes, Route } from 'react-router';
import App from './App';
import Auth from './auth/Index';

function Router() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/*" element={<App />} />
                <Route path="/login" element={<Auth type={'login'}/>} />
                <Route path="/register" element={<Auth type={'register'}/>} />
            </Routes>
        </BrowserRouter>
    );
}

export default Router;
