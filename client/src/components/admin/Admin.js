import React from 'react';
import { Link } from 'react-router-dom';
import './Admin.css'; // Import the CSS file

const Admin = () => {
  return (
    <div className="admin-container">
      <h2>Admin Panel</h2>
      <nav>
        <div className="category">
          <h3>Услуги</h3>
          <ul>
            <li><Link to="/admin/services/allServices">Все услуги</Link></li>
            <li><Link to="/admin/services/addService">Добавить услугу</Link></li>
            <li><Link to="/admin/services/redactService">Редактировать услугу</Link></li>
            <li><Link to="/admin/services/deleteService">Удалить услугу</Link></li>
          </ul>
        </div>
        <div className="category">
          <h3>Работники</h3>
          <ul>
            <li><Link to="/admin/employees/allEmployees">Все работники</Link></li>
            <li><Link to="/admin/employees/addEmployee">Добавить работника</Link></li>
            <li><Link to="/admin/employees/redactEmployee">Редактировать работника</Link></li>
            <li><Link to="/admin/employees/deleteEmployee">Удалить работника</Link></li>
          </ul>
        </div>
        <div className="category">
          <h3>Записи</h3>
          <ul>
            <li><Link to="/admin/appointments/allAppointments">Все записи</Link></li>
            <li><Link to="/admin/appointments/deleteAppointments">Удалить запись</Link></li>
          </ul>
        </div>
        <div className="category">
          <h3>Расписание</h3>
          <ul>
            <li><Link to="/admin/schedule/AllSchedules">Всё расписание</Link></li>
            <li><Link to="/admin/schedule/uploadNewSchedule">Загрузить новое расписание</Link></li>
            <li><Link to="/admin/schedule/redactSchedule">Редактировать расписание</Link></li>
            <li><Link to="/admin/schedule/delete-schedule">Удалить расписание</Link></li>
          </ul>
        </div>
      </nav>
    </div>
  );
};

export default Admin;
