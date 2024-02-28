document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.unmaskBtn').forEach(button => {
        button.addEventListener('click', function() {
            const token = this.getAttribute('data-token');
            const dataItem = this.closest('.dataItem');

            // Target the specific placeholders for update
            const cardNumberPlaceholder = dataItem.querySelector(`#cardNumber${token}`);
            const expiryDatePlaceholder = dataItem.querySelector(`#expiryDate${token}`);

            fetch('/unmask', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ token: token })
            })
            .then(response => response.json())
            .then(data => {
                console.log(data);
                cardNumberPlaceholder.innerHTML = `Card Number: ${data.card.card_number}`;
                expiryDatePlaceholder.innerHTML = `Expiry Date: ${data.card.expiry_date}`;

                button.style.display = 'none';
            })
            .catch(error => {
                console.error('Error unmasking card:', error);
            });
        });
    });
});
