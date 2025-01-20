const tg = window.Telegram.WebApp;
tg.expand();

const savePasswordBtn = document.getElementById("savePassword");
const viewPasswordsBtn = document.getElementById("viewPasswords");
const searchPasswordsBtn = document.getElementById("searchPasswords");
const formContainer = document.getElementById("formContainer");
const passwordsContainer = document.getElementById("passwordsContainer");

const BASE_URL = "http://3.79.247.241:8080/api";
const userID = new URLSearchParams(window.location.search).get("user_id");

savePasswordBtn.addEventListener("click", () => {
    formContainer.innerHTML = `
        <form id="saveForm">
            <h3>Yangi parol qo‘shish</h3>
            <input type="text" id="site" placeholder="Sayt nomi" required />
            <input type="text" id="password" placeholder="Parol" required />
            <button type="submit">Saqlash</button>
        </form>
    `;
    document.getElementById("saveForm").addEventListener("submit", async (e) => {
        e.preventDefault();
        const site = document.getElementById("site").value;
        const password = document.getElementById("password").value;

        try {
            const response = await fetch(`${BASE_URL}/post_password`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ user_id: userID, site, password }),
            });
            const result = await response.json();
            alert(result.message || "Parol muvaffaqiyatli saqlandi!");
        } catch (error) {
            alert("Xatolik yuz berdi: " + error.message);
        }
    });
});

viewPasswordsBtn.addEventListener("click", async () => {
    try {
        const response = await fetch(`${BASE_URL}/password/${userID}`);
        const result = await response.json();
        passwordsContainer.innerHTML = `<h3>Sizning parollaringiz</h3>`;
        result.data.forEach(password => {
            passwordsContainer.innerHTML += `<p><strong>${password.site}:</strong> ${password.password}</p>`;
        });
    } catch (error) {
        alert("Xatolik yuz berdi: " + error.message);
    }
});

searchPasswordsBtn.addEventListener("click", () => {
    formContainer.innerHTML = `
        <form id="searchForm">
            <h3>Parol qidirish</h3>
            <input type="text" id="searchSite" placeholder="Sayt nomi" required />
            <button type="submit">Qidirish</button>
        </form>
    `;
    document.getElementById("searchForm").addEventListener("submit", async (e) => {
        e.preventDefault();
        const site = document.getElementById("searchSite").value;

        try {
            const response = await fetch(`${BASE_URL}/password?userID=${userID}&site=${site}`);
            const result = await response.json();
            passwordsContainer.innerHTML = `<h3>Qidiruv natijalari</h3>`;
            result.data.forEach(password => {
                passwordsContainer.innerHTML += `<p><strong>${password.site}:</strong> ${password.password}</p>`;
            });
        } catch (error) {
            alert("Xatolik yuz berdi: " + error.message);
        }
    });
});