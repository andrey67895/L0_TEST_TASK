// Кнопка переключения темы
const themeToggleBtn = document.getElementById('themeToggleBtn');
themeToggleBtn.addEventListener('click', () => {
    document.body.classList.toggle('dark');
    if(document.body.classList.contains('dark')) {
        themeToggleBtn.textContent = 'Светлая тема';
    } else {
        themeToggleBtn.textContent = 'Темная тема';
    }
});

// Отправка формы
document.getElementById('orderForm').addEventListener('submit', async function(e){
    e.preventDefault();
    const uid = document.getElementById('uid').value.trim();
    if(!uid) return;

    try {
        const response = await fetch(`/order/${uid}`);
        if(response.ok) {
            const data = await response.json();
            document.getElementById('orderJson').textContent = JSON.stringify(data, null, 2);
            document.getElementById('orderInfo').style.display = 'block';
        } else {
            showError(`Ошибка ${response.status}: ${response.statusText}`);
        }
    } catch (err) {
        showError(`Сетевая ошибка: ${err.message}`);
    }
});

// Функция показа ошибки
function showError(message) {
    const popup = document.createElement('div');
    popup.className = 'error-popup';
    popup.textContent = message;
    document.body.appendChild(popup);
    setTimeout(() => {
        popup.style.opacity = '0';
        setTimeout(() => popup.remove(), 500);
    }, 5000);
}
