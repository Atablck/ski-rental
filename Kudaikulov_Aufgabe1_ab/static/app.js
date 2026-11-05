document.addEventListener("DOMContentLoaded", () => {
    loadDashboardStats()
    loadPersons()
    loadRentalOptions()
    loadInventory()
    loadActiveRentals()
});

function switchView(viewId) {
    document.querySelectorAll('.view').forEach(el => el.classList.remove('active'))
    document.getElementById(viewId).classList.add('active')
}

function loadDashboardStats() {
    fetch("/api/stats")
        .then(response => {
            if (!response.ok) throw new Error("Netzwerkantwort war nicht ok");
            return response.json();
        })
        .then(data => {
            document.getElementById("stat-available").innerText = data.available_items
            document.getElementById("stat-rented").innerText = data.rented_items
            document.getElementById("stat-total").innerText = data.total_items
            document.getElementById("stat-persons").innerText = data.total_persons
        })
        .catch(error => console.error("Fehler beim Laden der Statistik:", error))
}

function loadPersons() {
    fetch("/api/persons/list")
        .then(response => response.json())
        .then(data => {
            const tbody = document.getElementById("persons-list");
            tbody.innerHTML = ""

            if (!data || data.length === 0) {
                const tr = document.createElement("tr");
                tr.innerHTML = `<td colspan='4'>Keine Personen registriert</td>`;
                tbody.appendChild(tr);
                return
            }

            data.forEach(person => {
                const tr = document.createElement("tr");
                const ausweisLink = person.AusweisPfad ? `<br><a href="${person.AusweisPfad}" target="_blank" style="font-size: 12px; color: #3498db; text-decoration: underline;">Dokument ansehen</a>` : "";
                
                tr.innerHTML = `
                    <td>${person.ID}</td>
                    <td>${person.FirstName}</td>
                    <td>${person.LastName}</td>
                    <td>${person.AusweisNr} (${person.AusweisTyp}) ${ausweisLink}</td> 
                `;
                tbody.appendChild(tr)
            });
        })
        .catch(error => console.error("Fehler beim Laden der Personen:", error))
}

document.getElementById("person-form").addEventListener("submit", function(event) {
    event.preventDefault();

    let formData = new FormData(this)

    fetch("/api/persons", {
        method: "POST",
        body: formData
    })
    .then(response => {
        if (response.ok) {
            alert("Person erfolgreich registriert!")
            this.reset()
            loadDashboardStats()
            loadPersons()
            loadRentalOptions()
        } else {
            alert("Fehler beim Speichern")
        }
    })
    .catch(error => console.error("Fehler:", error))
})

function loadRentalOptions() {
    fetch("/api/persons/list").then(res => res.json()).then(data => {
        const select = document.getElementById("rental-person-select")
        select.innerHTML = ""
        if (data) {
            data.forEach(p => {
                select.innerHTML += `<option value="${p.ID}">${p.FirstName} ${p.LastName} (${p.AusweisNr})</option>`
            })
        }
    })

    fetch("/api/items/available").then(res => res.json()).then(data => {
        const select = document.getElementById("rental-item-select")
        select.innerHTML = ""
        if (data && data.length > 0) {
            data.forEach(item => {
                select.innerHTML += `<option value="${item.ID}">${item.Name} (${item.PricePerHour}€/h)</option>`
            })
        } else {
            select.innerHTML = `<option value="">Keine Ski verfügbar</option>`
        }
    })
}

document.getElementById("rental-form").addEventListener("submit", function(event) {
    event.preventDefault()

    let formData = new URLSearchParams(new FormData(this))

    fetch("/api/rentals", {
        method: "POST",
        body: formData,
        headers: { "Content-Type": "application/x-www-form-urlencoded" }
    })
    .then(response => {
        if (response.ok) {
            alert("Ausleihe erfolgreich!")
            this.reset();
            loadDashboardStats()
            loadRentalOptions()
            loadActiveRentals()
        } else {
            alert("Fehler bei der Ausleihe")
        }
    })
    .catch(error => console.error("Fehler:", error))
})

function loadInventory() {
    fetch("/api/items/all")
        .then(res => res.json())
        .then(data => {
            const tbody = document.getElementById("inventory-list")
            tbody.innerHTML = ""
            if (!data) return
            
            data.forEach(item => {
                const tr = document.createElement("tr")
                tr.innerHTML = `
                    <td>${item.id}</td>
                    <td><strong>${item.type}</strong></td>
                    <td>${item.name}</td>
                    <td>${item.price} €</td>
                `
                tbody.appendChild(tr)
            })
        })
}

function loadActiveRentals() {
    fetch("/api/rentals/active")
        .then(res => res.json())
        .then(data => {
            const tbody = document.getElementById("active-rentals-list")
            tbody.innerHTML = ""
            
            if (!data || data.length === 0) {
                tbody.innerHTML = "<tr><td colspan='4'>Keine aktiven Ausleihen</td></tr>"
                return
            }

            data.forEach(rent => {
                const tr = document.createElement("tr")
                tr.innerHTML = `
                    <td>${rent.person_name}</td>
                    <td>${rent.item_name}</td>
                    <td>${rent.hours}</td>
                    <td>
                        <button class="btn btn-orange btn-small" onclick="returnItem('${rent.rental_id}')">Rückgabe</button>
                    </td>
                `
                tbody.appendChild(tr)
            })
        })
}

function returnItem(rentalId) {
    if (!confirm("Ausleihe wirklich beenden?")) return

    let formData = new URLSearchParams()
    formData.append("rental_id", rentalId)

    fetch("/api/rentals/return", {
        method: "POST",
        body: formData,
        headers: { "Content-Type": "application/x-www-form-urlencoded" }
    })
    .then(res => res.json())
    .then(data => {
        alert("Rückgabe erfolgreich! Zu zahlen: " + data.price + " €")
        loadDashboardStats()
        loadRentalOptions()
        loadActiveRentals()
    })
    .catch(err => console.error("Fehler bei Rückgabe:", err))
}