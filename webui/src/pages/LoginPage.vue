<template>
  <div class="login-container">
    <div class="login-box">
      <h2>Login</h2>
      <form @submit.prevent="login">
        <input v-model="username" type="text" placeholder="Enter your username" required />
        <button type="submit">Login</button>
      </form>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      username: ''
    };
  },
  methods: {
    login() {
      if (!this.username || this.username.length < 2) {
        alert('Username must be at least 2 characters long');
        return;
      }

      axios.post('http://localhost:8080/login', { username: this.username })
        .then(response => {
          localStorage.setItem('access_token', response.data.access_token);
          localStorage.setItem('user_id', JSON.stringify(response.data.user_id));
          this.$router.push('/chats');
        })
        .catch(error => {
          console.error(error);
          alert('Invalid username');
        });
    }
  }
};
</script>

<style scoped>
:root {
  --primary-color: #42b883; 
  --secondary-color: #35495e; 
  --background-color: #ecf0f1; 
}

html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.login-container {
  padding: 0;
  margin: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  background: linear-gradient(135deg, #42b883, #35495e);
}

.login-box {
  background: white;
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.2);
  text-align: center;
  width: 350px;
  border-top: 5px solid var(--primary-color);
  transition: transform 0.3s ease;
}

.login-box:hover {
  transform: scale(1.02);
}

h2 {
  margin-bottom: 20px;
  color: var(--secondary-color);
  font-size: 22px;
  font-weight: bold;
}

input {
  width: 90%;
  padding: 12px;
  margin: 10px auto;
  display: block;
  border: 2px solid var(--secondary-color);
  border-radius: 5px;
  font-size: 16px;
  text-align: center;
  transition: border-color 0.3s ease;
}

input:focus {
  border-color: var(--primary-color);
  outline: none;
  box-shadow: 0 0 5px rgba(66, 184, 131, 0.5);
}

button {
  width: 100%;
  padding: 12px;
  background: #4bc599;
  color: white;
  font-weight: bold;
  border: none;
  border-radius: 5px;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.3s ease, transform 0.2s ease;
}

button:hover {
  background: #368c6c;
}

button:active {
  transform: scale(0.98);
}
</style>
