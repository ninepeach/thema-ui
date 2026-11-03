// thema-ui.js — small progressive enhancements. Every component works
// without it; this only adds:
//   [data-dialog-open="id"]  open a <dialog>
//   [data-copy="id"]         copy an element's text, show data-copied-label
//   details[data-menu]       close on outside click / Escape / when another opens
//   form[data-lock]          disable submit buttons while submitting (.is-loading)
//   input[type=file]         show the chosen file name in [data-file-name]
//   [data-scroll-progress]   set --progress (0–1) to how far the page is scrolled
(() => {
	"use strict";

	document.addEventListener("click", (event) => {
		const opener = event.target.closest("[data-dialog-open]");
		if (opener) {
			const dialog = document.getElementById(opener.dataset.dialogOpen);
			if (dialog && typeof dialog.showModal === "function") {
				event.preventDefault();
				dialog.showModal();
			}
			return;
		}

		const copier = event.target.closest("[data-copy]");
		if (copier) {
			const source = document.getElementById(copier.dataset.copy);
			const label = copier.querySelector("span");
			if (!source || !navigator.clipboard) return;
			navigator.clipboard.writeText(source.textContent.trim()).then(() => {
				if (!label) return;
				const original = label.textContent;
				label.textContent = copier.dataset.copiedLabel || original;
				copier.classList.add("is-done");
				setTimeout(() => {
					label.textContent = original;
					copier.classList.remove("is-done");
				}, 1600);
			});
			return;
		}

		document.querySelectorAll("details[data-menu][open]").forEach((menu) => {
			if (!menu.contains(event.target)) menu.removeAttribute("open");
		});
	});

	document.addEventListener("toggle", (event) => {
		const menu = event.target;
		if (!(menu instanceof HTMLDetailsElement) || !menu.open || !menu.hasAttribute("data-menu")) return;
		document.querySelectorAll("details[data-menu][open]").forEach((other) => {
			if (other !== menu) other.removeAttribute("open");
		});
	}, true);

	document.addEventListener("keydown", (event) => {
		if (event.key !== "Escape") return;
		document.querySelectorAll("details[data-menu][open]").forEach((menu) => {
			menu.removeAttribute("open");
			menu.querySelector("summary")?.focus();
		});
	});

	document.addEventListener("submit", (event) => {
		const form = event.target;
		if (!form.matches("form[data-lock]")) return;
		form.querySelectorAll('button[type="submit"]').forEach((button) => {
			button.disabled = true;
			button.classList.add("is-loading");
		});
		form.setAttribute("aria-busy", "true");
	});

	document.addEventListener("change", (event) => {
		const input = event.target;
		if (!input.matches('input[type="file"]')) return;
		const output = input.closest(".ui-file")?.querySelector("[data-file-name]");
		if (output) output.textContent = Array.from(input.files).map((f) => f.name).join(", ");
	});

	const bars = document.querySelectorAll("[data-scroll-progress]");
	if (bars.length) {
		let pending = false;
		const update = () => {
			pending = false;
			const root = document.documentElement;
			const total = root.scrollHeight - root.clientHeight;
			const progress = total > 0 ? Math.min(1, Math.max(0, window.scrollY / total)) : 0;
			bars.forEach((bar) => bar.style.setProperty("--progress", progress.toFixed(4)));
		};
		const schedule = () => {
			if (!pending) {
				pending = true;
				requestAnimationFrame(update);
			}
		};
		window.addEventListener("scroll", schedule, { passive: true });
		window.addEventListener("resize", schedule, { passive: true });
		update();
	}
})();
