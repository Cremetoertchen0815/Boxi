const sidebar = document.getElementById('sidebar');
const menuBtn = document.getElementById('menuBtn');
const overlay = document.getElementById('overlay');

function openSidebar() {
  sidebar.classList.add('open');
  overlay.classList.add('show');
}

function closeSidebar() {
  sidebar.classList.remove('open');
  overlay.classList.remove('show');
}

menuBtn.addEventListener('click', openSidebar);
overlay.addEventListener('click', closeSidebar);

// Swipe detection with Hammer.js
const hammer = new Hammer(document.body);
hammer.on('swiperight', (ev) => {
  if (!sidebar.classList.contains('open') && ev.deltaX > 30) {
    openSidebar();
  }
});
hammer.on('swipeleft', (ev) => {
  if (sidebar.classList.contains('open') && ev.deltaX < -30) {
    closeSidebar();
  }
});
