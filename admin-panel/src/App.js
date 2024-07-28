import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './components/Login';
import Admin from './components/admin/Admin';
import AllServices from './components/admin/Services/AllServices';
import AddService from './components/admin/Services/AddService';
import RedactService from './components/admin/Services/RedactService';
import DeleteService from './components/admin/Services/DeleteService';
import AllEmployees from './components/admin/Employees/AllEmployees';
import AddEmployee from './components/admin/Employees/AddEmployee';
import RedactEmployee from './components/admin/Employees/RedactEmployee';
import DeleteEmployee from './components/admin/Employees/DeleteEmployee';
import AllAppointments from './components/admin/Appointments/AllAppointments';
import DeleteAppointments from './components/admin/Appointments/DeleteAppointments';
import AllSchedules from './components/admin/Schedule/AllSchedules';
import UploadNewSchedule from './components/admin/Schedule/UploadNewSchedule';
import RedactSchedule from './components/admin/Schedule/RedactSchedule';
import DeleteSchedule from './components/admin/Schedule/DeleteSchedule';
import PrivateRoute from './components/PrivateRoute';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/admin" element={<PrivateRoute />}>
          <Route path="/admin" element={<Admin />} />
          <Route path="/admin/services/allServices" element={<AllServices />} />
          <Route path="/admin/services/addService" element={<AddService />} />
          <Route path="/admin/services/redactService" element={<RedactService />} />
          <Route path="/admin/services/deleteService" element={<DeleteService />} />
          <Route path="/admin/employees/allEmployees" element={<AllEmployees />} />
          <Route path="/admin/employees/addEmployee" element={<AddEmployee />} />
          <Route path="/admin/employees/redactEmployee" element={<RedactEmployee />} />
          <Route path="/admin/employees/deleteEmployee" element={<DeleteEmployee />} />
          <Route path="/admin/appointments/allAppointments" element={<AllAppointments />} />
          <Route path="/admin/appointments/DeleteAppointments" element={<DeleteAppointments />} />
          <Route path="/admin/schedule/allSchedules" element={<AllSchedules />} />
          <Route path="/admin/schedule/uploadNewSchedule" element={<UploadNewSchedule />} />
          <Route path="/admin/schedule/redactSchedule" element={<RedactSchedule />} />
          <Route path="/admin/schedule/delete-schedule" element={<DeleteSchedule />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
