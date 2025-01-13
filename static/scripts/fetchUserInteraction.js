function loadUserState() {
    const loginForm = document.getElementById('loginForm');
    const logoutForm = document.getElementById('logoutForm');
    const RegisterForm = document.getElementById('RegisterForm');
    const CreateForm = document.getElementById('postForm');

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
                this.reset();
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
                sendWebSocketMessage({type: 'logout'});
                navigateToPage('Login')
            }
        })
        .catch(error => {
            console.error('Error:', error);
            navigateToPage('Error');
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
                this.reset();
                navigateToPage('Home')
            } else {
                document.getElementById('register_err_field').innerHTML = result;
            }
        })
        .catch(error => {
            console.error('Error:', error);
            navigateToPage('Error');
        });
        
    });

    CreateForm.addEventListener('submit', function(event) {
        event.preventDefault();
    
        const formData = new FormData(this);        
        fetch('/Data-CreatPost', {
            method: 'POST',
            body: formData,
        })
        .then(response => response.text())
        .then(result => {
            if (result === 'Post Created Successfully') {
                const menuItems = document.querySelectorAll('.menu-item');
                menuItems.forEach(i => i.classList.remove('color'));    
                navigateToPage('Home')
                this.reset();
            } else {
                document.getElementById('CreatePost_err_field').innerHTML = result;
            }
        })
        .catch(error => {
            console.error('Error:', error);
            navigateToPage('Error');
        });
    });
}