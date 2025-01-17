const tg = window.Telegram.WebApp;
tg.expand(); // Telegram Web App-ni kengaytirish

const savePasswordBtn = document.getElementById("savePassword");
const viewPasswordsBtn = document.getElementById("viewPasswords");
const searchPasswordsBtn = document.getElementById("searchPasswords");
const formContainer = document.getElementById("formContainer");
const passwordsContainer = document.getElementById("passwordsContainer");

const BASE_URL = "http://3.79.247.241:8080/api/swagger/password";
const userID = tg.initDataUnsafe?.user?.id || "demo-user-id"; 

function getByName(site) {
    const url = `${BASE_URL}/search?userID=${encodeURIComponent(userID)}&site=${encodeURIComponent(site)}`;
    fetch(url)
        .then(response => {
            if (!response.ok) {
                throw new Error("No passwords found for the given criteria");
            }
            return response.json();
        })
        .then(data => {
            console.log("Fetched data:", data);
            const passwords = Array.isArray(data) ? data : [];
            passwordsContainer.innerHTML = `
                <h3>Qidiruv natijasi</h3>
                <ul>
                    ${passwords.map(p => `<li><strong>${p.site}</strong>: ${p.password}</li>`).join("")}
                </ul>
            `;
        })
        .catch(err => {
            console.error(err);
            passwordsContainer.innerHTML = `<p>${err.message}</p>`;
        });
}

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
            const response = await fetch(BASE_URL, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ site, password }),
            });
            const result = await response.json();
            alert(result.message || "Parol muvaffaqiyatli saqlandi!");
        } catch (error) {
            alert("Xatolik yuz berdi: " + error.message);
        }
    });
});

// Barcha parollarni ko‘rish
viewPasswordsBtn.addEventListener("click", async () => {
    try {
        const response = await fetch(`${BASE_URL}/all?userID=${encodeURIComponent(userID)}`);
        const passwords = await response.json();

        passwordsContainer.innerHTML = `
            <h3>Saqlangan parollar</h3>
            <ul>
                ${passwords
                    .map(
                        (password) =>
                            `<li><strong>${password.site}</strong>: ${password.password}</li>`
                    )
                    .join("")}
            </ul>
        `;
    } catch (error) {
        alert("Xatolik yuz berdi: " + error.message);
    }
});

// Parollarni qidirish
searchPasswordsBtn.addEventListener("click", () => {
    formContainer.innerHTML = `
        <form id="searchForm">
            <h3>Parollarni qidirish</h3>
            <input type="text" id="searchSite" placeholder="Sayt nomi" required />
            <button type="submit">Qidirish</button>
        </form>
    `;
    document.getElementById("searchForm").addEventListener("submit", async (e) => {
        e.preventDefault();
        const site = document.getElementById("searchSite").value;

        getByName(site);
        });
});
