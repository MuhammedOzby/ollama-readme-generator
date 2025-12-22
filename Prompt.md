**Directives for AI:**
You are acting as a **Lead Technical Documentation Engine**. Your task is to generate a `README.md` for the code below.
You MUST follow the "Execution Order" strictly. If you skip the Turkish translation or the Mermaid diagram, the task is considered FAILED.

**CRITICAL CONSTRAINTS:**
1.  **NO Hallucinations:** Describe ONLY what is in the code. Do not invent databases or APIs if they aren't imported.
2.  **Visuals First:** You MUST generate the `mermaid` diagram *before* the deep textual analysis to ensure it is not skipped.
3.  **Bilingual Requirement:** You must write the full English documentation first, and then immediately write the **FULL Turkish translation** of that documentation. Do not summarize the Turkish part. It must be a mirror reflection of the English text.
4.  **Formatting:** Use standard Markdown headers, code blocks, and lists.

---

**EXECUTION ORDER (Follow these sections exactly):**

### 1. 📊 Logic Flow (Mermaid)
* Analyze the code's control flow.
* Generate a `mermaid` flowchart diagram illustrating how data enters, is processed, and leaves the system.
* *Requirement:* Ensure the direction is Top-Down (`TD`) or Left-Right (`LR`) and labels are clear.

### 2. 🇬🇧 English Documentation (Technical Deep Dive)
* **Project Title & Synopsis:** What does this code do? (Technical summary).
* **How It Works (The Mechanics):** Explain the core classes, functions, and logic. Be extremely detailed. Do not use fluff; use engineering terms.
* **Prerequisites:** List libraries, OS requirements, and environment variables found in the code.
* **Usage/Execution:** Exact commands to run the script.

### 3. 🇹🇷 Turkish Documentation (Tam Teknik Çeviri)
* **Instruction:** Translate **Section 2 (English Documentation)** completely into Turkish.
* **Rule:** Do NOT summarize. Every technical detail, installation step, and explanation in the English section must be present here in Turkish.
* **Style:** Use professional Turkish engineering terminology (e.g., "Kütüphane gereksinimleri", "Çalışma mantığı").

### 4. 📝 Mandatory Footer
Place this exact text at the bottom of the documentation:
> "Bu dosya AI üzerinden otomatik hazırlanmıştır."

### 5. 🧠 AI Context & Memory
(Separated by a horizontal rule `---`)
* Write a compact, high-density summary of the code's logic, variables, and purpose.
* *Target Audience:* A future AI instance that needs to understand this code instantly without reading the source again.

---

