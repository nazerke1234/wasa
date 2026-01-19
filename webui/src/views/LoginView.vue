<script>
export default {
  name: "LoginView",
  data() {
    
    localStorage.clear();
    return {
      errorMessage: null, 
      userName: "",    
      userProfile: {      
        id: "",
        name: "",
      },
      isSubmitting: false, 
    };
  },
  methods: {
    
    async loadDefaultPhoto() {
      try {
        
        const response = await fetch('/nopfp.jpg'); 
        
        if (!response.ok) {
           throw new Error(`Failed to fetch default photo: ${response.statusText}`);
        }
        
        const blob = await response.blob();   
        const reader = new FileReader();   

        return new Promise((resolve, reject) => {
          reader.onload = () => {                     
            
            const base64String = reader.result.toString().split(',')[1];   
            resolve(base64String);
          };
          reader.onerror = reject;
          reader.readAsDataURL(blob);
        });
      } catch (error) {
        console.error('Error loading default photo:', error);
        return null; 
      }
    },

  
    async doLogin() {
      
      if (this.userName.trim() === "") {   
        this.errorMessage = "Name cannot be empty.";
        return;
      }
      
      this.errorMessage = null;
      this.isSubmitting = true; 

      try {
        const photoData = await this.loadDefaultPhoto();
        const response = await this.$axios.post("/session", {
          name: this.userName,
          photo: photoData,
        }, {
          headers: {
            'Content-Type': 'application/json' 
          }
        });
        
        const identifier = response.data.identifier;
        if (!identifier) {
          throw new Error("Unexpected server response. Missing 'identifier'.");
        }
        
        this.userProfile.id = identifier;
        this.userProfile.name = this.userName;
        
        
        localStorage.setItem("token", this.userProfile.id);
        localStorage.setItem("name", this.userProfile.name);
  

        this.$router.push({ path: "/home" });  

      } catch (e) {
        
        console.error("Login attempt failed:", e);  
        if (e.response) {
          if (e.response.status === 400) {
            this.errorMessage = "Input error: Please ensure your name meets all requirements.";
          } else if (e.response.status === 500) {
            this.errorMessage = "Server error: An internal error occurred. Please try again later.";
          }
        } else {
          this.errorMessage = "Network error or unexpected client issue.";
        }
      } finally {
        this.isSubmitting = false; 
      }
    },
  },
};
</script>

<template>
  <div class="login-page-wrapper">
    <div class="login-card">
      <h1 class="card-title">WASAText Messenger</h1>
      <p class="card-subtitle">Insert your name to log in or create a new account.</p>

      <div class="input-action-group">
        <input
          type="text"
          id="name"
          v-model="userName"
          class="app-input-field"
          placeholder="Your preferred name."
          @keyup.enter="doLogin"
        />
        <button 
          class="btn-primary login-button" 
          type="button" 
          @click="doLogin"
          :disabled="isSubmitting || userName.trim() === ''"
        >
          {{ isSubmitting ? 'Logging in...' : 'Start Messaging' }}
        </button>
      </div>

      <ErrorMsg v-if="errorMessage" :msg="errorMessage"></ErrorMsg>
    </div>
  </div>
</template>

<style scoped>




.login-page-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #ecf0f1; 
}

.login-card {
  max-width: 450px;
  width: 90%;
  padding: 30px;
  border-radius: 12px;
  background-color: white;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.card-title {
  font-size: 2.2rem;
  font-weight: 700;
  margin-bottom: 5px;
  color:  #4a90e2;
}

.card-subtitle {
  font-size: 1rem;
  color: #7f8c8d;
  margin-bottom: 30px;
}

.input-action-group {
  display: flex;
  flex-direction: column;
  gap: 15px;
  align-items: center;
}


.app-input-field {
  padding: 12px;
  width: 100%;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.2s;
}
.app-input-field:focus {
  border-color: #4a90e2;
  outline: none;
  box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.2);
}


.btn-primary {
  padding: 12px 25px;
  width: 100%;
  background-color:  #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1.1rem;
  font-weight: 600;
  transition: background-color 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background-color: #3a74b6;
}

.btn-primary:disabled {
    background-color: #ccc;
    cursor: not-allowed;
    opacity: 0.8;
}


@media (max-width: 500px) { 
  .login-card {
    margin: 10% auto;
  }
}
</style>
