document.addEventListener('DOMContentLoaded', function () {
    const expiryInput = document.getElementById('expiry_date');
    expiryInput.addEventListener('input', function(e) {
        let inputVal = expiryInput.value.replace(/\D/g, '');

        inputVal = inputVal.slice(0, 4);

        if (inputVal.length > 2) {
            inputVal = `${inputVal.slice(0, 2)} / ${inputVal.slice(2)}`;
        }

        expiryInput.value = inputVal;
    });

    const tokeniseForm = document.getElementById('tokeniseForm');
    if (tokeniseForm) {
        tokeniseForm.addEventListener('submit', function(event) {
            event.preventDefault();

            let expiryDateValue = tokeniseForm.expiry_date.value;
            expiryDateValue = expiryDateValue.split(" / ").join("");

            const requestData = {
                request_id: 'req-' + Math.random().toString(36).substring(2, 9),
                card: {
                    cardholder_name: tokeniseForm.cardholder_name.value,
                    card_number: tokeniseForm.card_number.value,
                    expiry_date: expiryDateValue,
                }
            };

            console.log(requestData)

            fetch('/tokenise', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData),
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    document.getElementById('responsePlaceholder').innerHTML = `<p>Some error occured: ${data.error}</p`;
                } else {
                    document.getElementById('responsePlaceholder').innerHTML = `<p>Token: ${data.token}</p>`;
                }
            })
            tokeniseForm.reset()
        });
    }

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
