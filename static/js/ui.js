export const UI = {
    showAuth: () => {
        document.getElementById("auth-screen").style.display = "block";
        document.getElementById("app-screen").classList.remove("visible");
    },
    
    showApp: (email) => {
        document.getElementById("auth-screen").style.display = "none";
        document.getElementById("app-screen").classList.add("visible");
        document.getElementById("user-email-display").textContent = email || "";
    },

    switchTab: (btn, mode) => {
        document.querySelectorAll(".auth-tab").forEach((t) => t.classList.remove("active"));
        btn.classList.add("active");
        document.getElementById("btn-auth-submit").textContent = mode;
        document.getElementById("auth-status").textContent = "";
    },

    setAuthStatus: (msg, type) => {
        const el = document.getElementById("auth-status");
        el.textContent = msg;
        el.className = `auth-status ${type}`;
    },

    togglePanel: (id) => {
        ["seed-panel", "import-panel"].forEach((p) => {
            if (p === id) document.getElementById(p).classList.toggle("open");
            else document.getElementById(p).classList.remove("open");
        });
    },

    toggleClear: () => {
        document.getElementById("confirm-box").classList.toggle("open");
    },

    toggleNotes: () => {
        const box = document.getElementById("notes-box");
        box.style.display = box.style.display === "none" ? "block" : "none";
    },

    showEmpty: (msg) => {
        document.getElementById("q-content").style.display = "none";
        const el = document.getElementById("empty");
        el.style.display = "flex";
        el.querySelector(".empty-text").textContent = msg || "no questions found";
    },

    showQuestion: (q) => {
        document.getElementById("empty").style.display = "none";
        document.getElementById("q-content").style.display = "block";
        document.getElementById("notes-box").style.display = "none";
        document.getElementById("q-title").textContent = q.title;
        
        const descEl = document.getElementById("q-desc");
        if (q.description) {
            descEl.textContent = q.description;
            descEl.style.display = "block";
        } else {
            descEl.style.display = "none";
        }

        document.getElementById("q-id").textContent = "#" + q.id;
        document.getElementById("notes-input").value = q.notes || "";
        document.getElementById("q-topic").textContent = q.topic;
        document.getElementById("q-title").textContent = q.title;
        document.getElementById("q-leetcode").href = `https://leetcode.com/problems/${q.slug}`;  
      
        const diffEl = document.getElementById("q-diff");
        diffEl.textContent = q.difficulty;
        diffEl.className = "diff-badge diff-" + q.difficulty;
        
        UI.updateActionButtons(q);
    },

    updateActionButtons: (q) => {
        document.getElementById("btn-solved").className = "btn-action" + (q.status === "solved" ? " active-solved" : "");
        document.getElementById("btn-skipped").className = "btn-action" + (q.status === "skipped" ? " active-skipped" : "");
        document.getElementById("btn-bookmark").className = "btn-action" + (q.bookmarked ? " active-bookmarked" : "");
        
        const badge = document.getElementById("q-status-badge");
        if (q.status) {
            badge.textContent = q.status;
            badge.className = "status-badge status-" + q.status;
            badge.style.display = "inline-block";
        } else {
            badge.style.display = "none";
        }
    },

    setSeedStatus: (msg, type) => {
        const el = document.getElementById("seed-status");
        el.textContent = msg;
        el.className = `status-msg ${type}`;
    },

    setImportStatus: (msg, type) => {
        const el = document.getElementById("import-status");
        el.textContent = msg;
        el.className = `status-msg ${type}`;
    },

    clearSeedForm: () => {
        document.getElementById("s-title").value = "";
        document.getElementById("s-topic").value = "";
        document.getElementById("s-diff").value = "";
        document.getElementById("s-slug").value = "";
    }
};