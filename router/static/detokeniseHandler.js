document.addEventListener('DOMContentLoaded', function () {
    const detokeniseForm = document.getElementById('detokeniseForm');
    if (detokeniseForm) {
        detokeniseForm.addEventListener('submit', function(event) {
            event.preventDefault();

            const requestData = {
                request_id: 'req-' + Math.random().toString(36).substring(2, 9),
                token: detokeniseForm.token.value,
            };

            fetch('/detokenise', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData),
            })
            .then(response => response.json())
            .then(data => {
                const content = `<p>Card Number: ${data.card.card_number}</p>
                                 <p>Expiry Date: ${data.card.expiry_date}</p>`;
                document.getElementById('responsePlaceholder').innerHTML = content;
            })
            .catch((error) => {
                console.error('Error:', error);
                document.getElementById('responsePlaceholder').innerHTML = `<p>Error: ${error.toString()}</p>`;
            });
        });
    }
});
