const API_BASE_URL = 'http://localhost:8080/api';

function showAlert(elementId, message, type = 'success') {
    const alert = document.getElementById(elementId);
    if (!alert) return;

    alert.textContent = message;
    alert.className = `alert alert-${type}`;
    alert.style.display = 'block';

    setTimeout(() => {
        alert.style.display = 'none';
    }, 5000);
}

async function saveData() {
    const key = document.getElementById('keyInput').value.trim();
    const value = document.getElementById('valueInput').value.trim();

    if (!key || !value) {
        showAlert('saveAlert', 'Please enter both key and value', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/data`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ key, value })
        });

        if (response.ok) {
            const result = await response.json();
            showAlert('saveAlert', `${result.status}: ${result.key}`);
            document.getElementById('keyInput').value = '';
            document.getElementById('valueInput').value = '';

            // Обновляем данные после сохранения
            getAllData();
        } else {
            const error = await response.json();
            showAlert('saveAlert', `${error.error}`, 'error');
        }
    } catch (error) {
        showAlert('saveAlert', `Network error: ${error.message}`, 'error');
    }
}

async function getAllData() {
    try {
        const response = await fetch(`${API_BASE_URL}/data`);

        if (response.ok) {
            const data = await response.json();
            const dataElement = document.getElementById('allData');

            if (!dataElement) return;

            if (Object.keys(data).length === 0) {
                dataElement.textContent = 'Database is empty';
            } else {
                dataElement.textContent = JSON.stringify(data, null, 2);
            }

            showAlert('dataAlert', 'Data loaded successfully!');
        } else {
            showAlert('dataAlert', 'Failed to load data', 'error');
        }
    } catch (error) {
        const dataElement = document.getElementById('allData');
        if (dataElement) {
            dataElement.textContent = 'Error: ' + error.message;
        }
        showAlert('dataAlert', 'Network error', 'error');
    }
}

async function deleteData() {
    const key = document.getElementById('deleteKeyInput').value.trim();

    if (!key) {
        showAlert('deleteAlert', 'Please enter a key to delete', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/data/${key}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            const result = await response.json();
            showAlert('deleteAlert', `${result.status}: ${result.key}`);
            document.getElementById('deleteKeyInput').value = '';

            // Обновляем данные после удаления
            getAllData();
        } else if (response.status === 404) {
            showAlert('deleteAlert', `Key "${key}" not found`, 'error');
        } else {
            const error = await response.json();
            showAlert('deleteAlert', `${error.error}`, 'error');
        }
    } catch (error) {
        showAlert('deleteAlert', `Network error: ${error.message}`, 'error');
    }
}

async function getStats() {
    try {
        const response = await fetch(`${API_BASE_URL}/stats`);

        if (response.ok) {
            const stats = await response.json();

            const totalRequestsElement = document.getElementById('totalRequests');
            const dbSizeElement = document.getElementById('dbSize');
            const statsDataElement = document.getElementById('statsData');

            if (totalRequestsElement) {
                totalRequestsElement.textContent = stats.total_requests;
            }

            if (dbSizeElement) {
                dbSizeElement.textContent = stats.database_size;
            }

            if (statsDataElement) {
                statsDataElement.textContent = JSON.stringify(stats, null, 2);
            }

            showAlert('statsAlert', 'Statistics loaded!');
        } else {
            showAlert('statsAlert', 'Failed to load statistics', 'error');
        }
    } catch (error) {
        showAlert('statsAlert', `Network error: ${error.message}`, 'error');
    }
}

document.addEventListener('DOMContentLoaded', function() {
    getAllData();

    document.getElementById('keyInput')?.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') saveData();
    });

    document.getElementById('valueInput')?.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') saveData();
    });

    document.getElementById('deleteKeyInput')?.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') deleteData();
    });
});