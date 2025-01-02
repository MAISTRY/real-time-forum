function loadUserState() {
    const loginForm = document.getElementById('loginForm');
    const logoutForm = document.getElementById('logoutForm');
    const RegisterForm = document.getElementById('RegisterForm');

    loginForm.addEventListener('submit', function(event) {
        event.preventDefault();
    
        const formData = new FormData(this);
        const data = new URLSearchParams(formData);
    
        fetch('/Data-userLogin', {
            method: 'POST',
            body: data,
        })
        .then(response => response.text())
        .then(result => {
            if (result === 'Login successful') {
                Authenticated();
                navigateToPage('Home')
            } else {
                document.getElementById('login_err_field').innerHTML = result;
            }
        })
        .catch(error => {
            console.error('Error:', error);
            navigateToPage('Error');
        });
    });

    logoutForm.onclick = (event) => {
        event.preventDefault(); 
        
        fetch('/Data-userLogout', {
            method: 'POST',
        })
        .then(response => response.text())
        .then(result => {
            if (result === 'Logout successful') {
                Authenticated();
                navigateToPage('Login')
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
    }

    RegisterForm.addEventListener('submit', function(event) {
        event.preventDefault();
    
        const formData = new FormData(this);
        const data = new URLSearchParams(formData);
    
        fetch('/Data-userRegister', {
            method: 'POST',
            body: data,
        })
        .then(response => response.text())
        .then(result => {
            if (result === 'Registration successful') {
                Authenticated();
                navigateToPage('Home')
            } else {
                document.getElementById('login_err_field').innerHTML = result;
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
}