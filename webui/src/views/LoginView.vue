<template>
  <div class="login-wrapper">
    <div class="login-box">
      <h2 class="welcome-text">Welcome to WASAText</h2>
      <div class="form-section">
        <input
          v-model="userName"
          type="text"
          placeholder="Enter your name"
          class="name-input"
        />
        <button @click="handleLogin" class="submit-btn">Log In</button>
      </div>
      <ErrorMsg v-if="error" :msg="error" />
    </div>
  </div>
</template>

<script>
import ErrorMsg from '../components/ErrorMsg.vue';

export default {
  components: { ErrorMsg },
  data() {
    localStorage.clear();
    return {
      userName: '',
      error: null,
      identity: {
        token: '',
        displayName: ''
      }
    };
  },
  methods: {
    async fetchDefaultImage() {
      try {
        const res = await fetch('/nopfp.jpg');
        const blob = await res.blob();
        return await this.toBase64(blob);
      } catch (e) {
        console.error('Failed to fetch default image:', e);
        return null;
      }
    },
    toBase64(blob) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onloadend = () => resolve(reader.result.split(',')[1]);
        reader.onerror = reject;
        reader.readAsDataURL(blob);
      });
    },
    async handleLogin() {
      if (!this.userName.trim()) {
        this.error = 'Name is required.';
        return;
      }

      try {
        const defaultPhoto = await this.fetchDefaultImage();
        const response = await this.$axios.post('/session', {
          name: this.userName,
          photo: defaultPhoto
        }, {
          headers: { 'Content-Type': 'application/json' }
        });

        const id = response.data.identifier;
        if (!id) {
          throw new Error('Missing identifier in response.');
        }

        this.identity.token = id;
        this.identity.displayName = this.userName;

        localStorage.setItem('token', id);
        localStorage.setItem('name', this.userName);

        this.$router.push('/home');
      } catch (err) {
        if (err.response?.status === 400) {
          this.error = 'Please fill in the name correctly.';
        } else if (err.response?.status === 500) {
          this.error = 'Server error occurred.';
        } else {
          this.error = err.message || 'Login failed.';
        }
      }
    }
  }
};
</script>

<style scoped>
.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 3rem;
}

.login-box {
  background: #f5f8fa;
  padding: 2rem;
  border-radius: 10px;
  max-width: 400px;
  width: 100%;
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.welcome-text {
  font-size: 22px;
  font-weight: bold;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.name-input {
  padding: 0.75rem;
  border: 1px solid #ccc;
  border-radius: 6px;
  font-size: 14px;
  width: 100%;
}

.submit-btn {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 0.75rem;
  border-radius: 6px;
  font-size: 15px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.submit-btn:hover {
  background-color: #0056b3;
}
</style>
