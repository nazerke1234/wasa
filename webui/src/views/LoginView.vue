<script>
export default {
  data() {
    localStorage.clear();
    return {
      errorNotification: null,
      usernameInput: "", 
      userProfile: {
        identifier: "",
        displayName: "",
      },
    };
  },
  methods: {
    async fetchDefaultAvatar() {
      try {
          const imageResponse = await fetch('/nopfp.jpg');
          const imageBlob = await imageResponse.blob();
          const fileReader = new FileReader();
          return new Promise((resolve, reject) => {
              fileReader.onload = () => {
                  const base64Data = fileReader.result.toString().split(',')[1];
                  resolve(base64Data);
              };
              fileReader.onerror = reject;
              fileReader.readAsDataURL(imageBlob);
          });
      } catch (fetchError) {
          console.error('Unable to load default avatar:', fetchError);
          return null;
      }
    },
    async processLogin() {
      if (this.usernameInput.trim() === "") {
        this.errorNotification = "Please enter your username.";
        return;
      }
      try {
        const avatarData = await this.fetchDefaultAvatar();
        const apiResponse = await this.$axios.post("/session", {
          name: this.usernameInput,
          photo: avatarData,
        }, {
          headers: {
            'Content-Type': 'application/json'
          }
        });
        if (apiResponse.data.identifier) {
          this.userProfile.identifier = apiResponse.data.identifier;
          this.userProfile.displayName = this.usernameInput; 
        } else {
          throw new Error("Invalid server response. Missing user identifier.");
        }
        localStorage.setItem("token", this.userProfile.identifier);
        localStorage.setItem("name", this.userProfile.displayName);
        this.$router.push({ path: "/home" });
      } catch (apiError) {
        if (apiError.response && apiError.response.status === 400) {
          this.errorNotification =
            "Please verify all fields and try again.";
        } else if (apiError.response && apiError.response.status === 500) {
          this.errorNotification =
            "Server error occurred. Please try again later.";
        } else {
          this.errorNotification = apiError.toString();
        }
      }
    },
  },
};
</script>

<template>
  <div class="auth-panel">
    <h1 class="app-welcome">Welcome to WASAText Messenger</h1>
    <div class="input-section">
      <input
        type="text"
        id="username-field"
        v-model="usernameInput"
        class="username-input"
        placeholder="Enter your username to access WASAText"
        @keyup.enter="processLogin"
      />
      <button class="auth-button" type="button" @click="processLogin">Enter</button>
    </div>
    <ErrorMsg v-if="errorNotification" :msg="errorNotification"></ErrorMsg>
  </div>
</template>

<style scoped>
.auth-panel {
  max-width: 420px;
  margin: 80px auto;
  text-align: center;
  padding: 30px;
  border: 1px solid #e1e5e9;
  border-radius: 12px;
  background-color: #ffffff;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.app-welcome {
  font-size: 26px;
  font-weight: 600;
  margin-bottom: 24px;
  color: #2c3e50;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.input-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.username-input {
  padding: 12px 16px;
  width: 100%;
  margin-bottom: 8px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
}

.username-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.auth-button {
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 600;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  width: 100%;
}

.auth-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.auth-button:active {
  transform: translateY(0);
}

@media (max-width: 480px) {
  .auth-panel {
    margin: 40px 20px;
    padding: 24px;
  }
  
  .app-welcome {
    font-size: 22px;
  }
}
</style>
