const sidebar = document.getElementById('sidebar');
const menuBtn = document.getElementById('menuBtn');
const overlay = document.getElementById('overlay');
const baseAddr = "http://192.168.4.1:8080/"

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

$(document).on('change blur', 'input[type="number"]', function () {
  const $input = $(this);
  const minAttr = $input.attr('min');
  const maxAttr = $input.attr('max');

  if (minAttr === undefined || maxAttr === undefined) return;

  const min = parseFloat(minAttr);
  const max = parseFloat(maxAttr);
  let val = parseFloat($input.val());

  // Only validate if it's a number
  if (!isNaN(val)) {
    if (min !== null && val < min) {
      $input.val(min);
    } else if (max !== null && val > max) {
      $input.val(max);
    }
  } else {
    $input.val(min);
  }
});