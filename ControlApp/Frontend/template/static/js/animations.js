
$('#animation-add-icon').on('click', e => {
   e.target.style.display = 'none';

   $('#animation-add-container')[0].classList.add('animation-add-container-active');
   $('#animation-add-dialog')[0].style.display = 'block';
});