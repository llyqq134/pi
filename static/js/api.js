const API_BASE = 'http://localhost:8080';

function getToken() {
  return localStorage.getItem('auth_token');
}

function getCurrentUser() {
  try {
    const user = JSON.parse(localStorage.getItem('current_user') || '{}');
    return user;
  } catch {
    return null;
  }
}

function isLoggedIn() {
  return !!getToken();
}

function requireAuth() {
  if (!isLoggedIn()) {
    window.location.href = '/';
    return false;
  }
  return true;
}

async function apiRequest(url, options = {}) {
  const token = getToken();
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  };
  if (token) {
    headers['Authorization'] = 'Bearer ' + token;
  }
  const res = await fetch(API_BASE + url, { ...options, headers });
  return res;
}

async function login(name, password) {
  const res = await fetch(API_BASE + '/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, password })
  });
  const data = await res.json();
  if (data.success && data.token) {
    localStorage.setItem('auth_token', data.token);
    localStorage.setItem('current_user', JSON.stringify(data.worker));
    return { success: true };
  }
  return { success: false, message: data.message || 'Ошибка входа' };
}

function logout() {
  localStorage.removeItem('auth_token');
  localStorage.removeItem('current_user');
  window.location.href = '/';
}

async function getWorkersByDepartment(department) {
  const res = await apiRequest('/listworkers/department/' + encodeURIComponent(department));
  if (!res.ok) throw new Error('Ошибка загрузки');
  return res.json();
}

async function addWorker(name, jobtitle, password) {
  const res = await apiRequest('/listworkers', {
    method: 'POST',
    body: JSON.stringify({ name, jobtitle, password })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || 'Ошибка');
  return data;
}

async function deleteWorker(id) {
  const res = await apiRequest('/listworkers/' + id, { method: 'DELETE' });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || 'Ошибка');
  return data;
}

async function getDepartments() {
  const res = await apiRequest('/departments/');
  if (!res.ok) throw new Error('Ошибка загрузки');
  return res.json();
}

async function addDepartment(name) {
  const res = await apiRequest('/departments/new', {
    method: 'POST',
    body: JSON.stringify({ name })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || 'Ошибка');
  return data;
}

async function deleteDepartment(name) {
  const res = await apiRequest('/departments/' + encodeURIComponent(name), {
    method: 'DELETE'
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || 'Ошибка');
  return data;
}

async function getEquipment() {
  const res = await fetch(API_BASE + '/equipment/list');
  if (!res.ok) throw new Error('Ошибка загрузки');
  return res.json();
}

async function addRecord(equipmentId, expectedReturnDate, status, departmentId, departmentName, workerId, workerName) {
  const body = {
    equipment_id: equipmentId,
    expected_return_date: expectedReturnDate,
    status: status || 'issued'
  };
  if (departmentId && departmentName && workerId && workerName) {
    body.department_id = departmentId;
    body.department_name = departmentName;
    body.worker_id = workerId;
    body.worker_name = workerName;
  }
  const res = await apiRequest('/record/add', {
    method: 'POST',
    body: JSON.stringify(body)
  });
  let data;
  try {
    data = await res.json();
  } catch (_) {
    throw new Error(res.ok ? 'Invalid response' : 'Server error');
  }
  if (!res.ok) throw new Error(data.message || 'Server error');
  return data;
}

async function downloadReport(startDate, endDate) {
  let url = API_BASE + '/record/export?';
  if (startDate) url += 'start_date=' + encodeURIComponent(startDate) + '&';
  if (endDate) url += 'end_date=' + encodeURIComponent(endDate);
  const token = getToken();
  const res = await fetch(url, {
    method: 'POST',
    headers: token ? { 'Authorization': 'Bearer ' + token } : {}
  });
  if (!res.ok) throw new Error('Ошибка формирования отчёта');
  const blob = await res.blob();
  const a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = 'records_' + (startDate || '') + '_' + (endDate || '') + '.pdf';
  a.click();
  URL.revokeObjectURL(a.href);
}
