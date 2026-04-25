console.log("main.js loaded");

import * as API from "./api.js";
import { UI } from "./ui.js";

// --- State ---
let currentQuestion = null;
let authMode = "login";
let selTopic = "";
let selDiff = "";

// --- Taglines ---
const taglines = [
  "because you always pick easy problems",
  "stop solving two sum for the 10th time",
  "fate decides your interview prep",
  "the algorithm chooses. you suffer.",
  "you can't skip what you don't pick",
  "no favorites. no excuses.",
  "your comfort zone called. we hung up.",
  "luck-based suffering, structured learning",
  "the universe knows you've been avoiding trees",
  "random problems, predictable excuses",
  "spin it. hate it. solve it.",
  "you've done two sum 7 times. we counted.",
  "dp is not that scary. spin and find out.",
  "FAANG doesn't care about your comfort zone.",
  "you said you'd do graphs tomorrow. tomorrow is now.",
  "linked lists don't care about your feelings",
  "you skipped hard mode again, didn't you",
  "the recruiter won't ask about two sum. we will.",
  "your leetcode streak is a lie. spin.",
  "somewhere a senior engineer is watching. spin faster.",
  "you memorized two sum. cute.",
  "no topic filters. face everything.",
  "graphs are just trees that went to therapy",
  "one spin. one problem. no excuses.",
  "your next interview is closer than you think"
]

function setRandomTagline() {
  console.log("setting tagline");
  const tagline = taglines[Math.floor(Math.random() * taglines.length)];
  document.getElementById("tagline").textContent = tagline;
  document.getElementById("tagline-app").textContent = tagline;
}

// --- Initialization ---
function init() {
    setRandomTagline();
    if (localStorage.getItem("token")) {
        UI.showApp(localStorage.getItem("userEmail"));
    } else {
        UI.showAuth();
    }
}

init();

function handleUnauthorized(res) {
    if (res.status === 401) {
        logout();
        return true;
    }
    return false;
}

function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("userEmail");
    currentQuestion = null;
    UI.showAuth();
}

// --- Auth Event Listeners ---
document.querySelectorAll(".auth-tab").forEach(btn => {
    btn.addEventListener("click", (e) => {
        authMode = e.target.dataset.mode;
        UI.switchTab(e.target, authMode);
    });
});

document.getElementById("btn-auth-submit").addEventListener("click", async () => {
    const email = document.getElementById("auth-email").value.trim();
    const password = document.getElementById("auth-password").value;

    if (!email || !password) {
        return UI.setAuthStatus("email and password required", "err");
    }

    try {
        const { res, data } = await API.loginOrRegister(authMode, email, password);
        if (!res.ok) throw new Error(data.error);

        if (authMode === "register") {
            UI.setAuthStatus("registered · please login", "ok");
            authMode = "login";
            UI.switchTab(document.getElementById("tab-login"), "login");
            return;
        }

        localStorage.setItem("token", data.token);
        localStorage.setItem("userEmail", email);
        UI.showApp(email);
    } catch (e) {
        UI.setAuthStatus("error · " + e.message, "err");
    }
});

document.getElementById("btn-logout").addEventListener("click", logout);

// --- Filter Event Listeners ---
document.querySelectorAll("#topics .pill").forEach((p) => {
    p.addEventListener("click", () => {
        document.querySelectorAll("#topics .pill").forEach((x) => x.classList.remove("active"));
        p.classList.add("active");
        selTopic = p.dataset.val;
    });
});

document.querySelectorAll("#diffs .pill").forEach((p) => {
    p.addEventListener("click", () => {
        document.querySelectorAll("#diffs .pill").forEach((x) => x.classList.remove("active"));
        p.classList.add("active");
        selDiff = p.dataset.val;
    });
});

// --- Action Event Listeners ---
document.getElementById("btn-toggle-seed").addEventListener("click", () => UI.togglePanel("seed-panel"));
document.getElementById("btn-toggle-import").addEventListener("click", () => UI.togglePanel("import-panel"));
document.getElementById("btn-toggle-clear").addEventListener("click", UI.toggleClear);
document.getElementById("btn-cancel-clear").addEventListener("click", UI.toggleClear);
document.getElementById("btn-notes-toggle").addEventListener("click", UI.toggleNotes);

document.getElementById("btn-spin").addEventListener("click", async () => {
    const btn = document.getElementById("btn-spin");
    const card = document.getElementById("card");
    btn.disabled = true;
    btn.querySelector("span").textContent = "...";
    card.classList.add("loading");

    try {
        const res = await API.fetchQuestion(selTopic, selDiff);
        if (handleUnauthorized(res)) return;
        
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || "not found");
        
        currentQuestion = data;
        UI.showQuestion(data);
    } catch (e) {
        currentQuestion = null;
        UI.showEmpty(e.message);
    } finally {
        btn.disabled = false;
        btn.querySelector("span").textContent = "spin";
        card.classList.remove("loading");
    }
});

document.getElementById("btn-solved").addEventListener("click", () => setStatus("solved"));
document.getElementById("btn-skipped").addEventListener("click", () => setStatus("skipped"));

async function setStatus(status) {
    if (!currentQuestion) return;
    const newStatus = currentQuestion.status === status ? null : status;
    try {
        const res = await API.updateQuestionStatus(currentQuestion.id, newStatus);
        if (handleUnauthorized(res)) return;
        currentQuestion.status = newStatus;
        UI.updateActionButtons(currentQuestion);
    } catch (e) {}
}

document.getElementById("btn-bookmark").addEventListener("click", async () => {
    if (!currentQuestion) return;
    const newVal = !currentQuestion.bookmarked;
    try {
        const res = await API.updateQuestionBookmark(currentQuestion.id, newVal);
        if (handleUnauthorized(res)) return;
        currentQuestion.bookmarked = newVal;
        UI.updateActionButtons(currentQuestion);
    } catch (e) {}
});

document.getElementById("btn-notes-save").addEventListener("click", async () => {
    if (!currentQuestion) return;
    const notes = document.getElementById("notes-input").value.trim() || null;
    const saveBtn = document.getElementById("btn-notes-save");
    
    saveBtn.textContent = "saving...";
    saveBtn.disabled = true;

    try {
        const res = await API.updateQuestionNotes(currentQuestion.id, notes);
        if (handleUnauthorized(res)) return;
        if (!res.ok) throw new Error("failed");
        
        currentQuestion.notes = notes;
        saveBtn.textContent = "saved!";
        saveBtn.style.color = "#5dcaa5";
        saveBtn.style.borderColor = "#5dcaa5";
    } catch (e) {
        saveBtn.textContent = "error!";
        saveBtn.style.color = "#f09595";
        saveBtn.style.borderColor = "#f09595";
    } finally {
        setTimeout(() => {
            saveBtn.textContent = "save notes";
            saveBtn.disabled = false;
            saveBtn.style.color = "";
            saveBtn.style.borderColor = "";
        }, 2000);
    }
});

// --- API Submissions ---
document.getElementById("btn-seed-submit").addEventListener("click", async () => {
    const title = document.getElementById("s-title").value.trim();
    const topic = document.getElementById("s-topic").value.trim();
    const difficulty = document.getElementById("s-diff").value;
    const slug = document.getElementById("s-slug").value.trim();
    
    if (!title || !topic || !difficulty || !slug) {
        return UI.setSeedStatus("all fields required", "err");
    }

    try {
        const res = await API.postQuestion({ title, topic, difficulty, slug });
        if (handleUnauthorized(res)) return;
        
        const data = await res.json();
        if (!res.ok) throw new Error(data.error);
        
        UI.setSeedStatus("added · " + data.slug, "ok");
        UI.clearSeedForm();
    } catch (e) {
        UI.setSeedStatus("error · " + e.message, "err");
    }
});

document.getElementById("btn-import-submit").addEventListener("click", async () => {
    const fileInput = document.getElementById("import-file");
    if (!fileInput.files.length) {
        return UI.setImportStatus("select a json file first", "err");
    }

    try {
        const text = await fileInput.files[0].text();
        const questions = JSON.parse(text);
        
        const res = await API.postImport(questions);
        if (handleUnauthorized(res)) return;
        
        const data = await res.json();
        if (!res.ok) throw new Error(data.error);
        
        UI.setImportStatus("imported · " + data.imported + " questions", "ok");
        fileInput.value = "";
    } catch (e) {
        UI.setImportStatus("error · invalid json or server error", "err");
    }
});

document.getElementById("btn-confirm-clear").addEventListener("click", async () => {
    try {
        const res = await API.deleteQuestions();
        if (handleUnauthorized(res)) return;
        
        const data = await res.json();
        if (!res.ok) throw new Error(data.error);
        
        UI.toggleClear();
        UI.showEmpty("spin to begin");
    } catch (e) {
        UI.toggleClear();
    }
});