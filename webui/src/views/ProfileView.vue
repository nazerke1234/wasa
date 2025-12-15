<template>
  <div class="profile-page">
    
    <header class="d-flex justify-content-between align-items-center pt-3 pb-2 mb-3 profile-topbar-custom">
      <h2 class="username-header">{{ userName }}, welcome to your profile</h2>
      
      <button type="button" class="btn-chat-tertiary" @click="logOut">
        Log Out
      </button>
    </header>

    <!-- Profile Card -->
    <div class="profile-card">
      <!-- Avatar -->
      <div class="avatar-wrapper">
        <img
          v-if="userPhoto"
          :src="userPhoto"
          class="avatar"
          alt="Profile photo"
        />
        <div v-else class="avatar placeholder">No photo</div>
      </div>

      <!-- Update photo -->
      <div class="form-group">
        <input type="file" @change="handlePhotoUpload" />
        <button
          class="primary-btn"
          @click="updatePhoto"
          :disabled="!newPhoto"
        >
          Update Photo
        </button>
      </div>

      <!-- User Info -->
      <div class="form-group">
        <input
            v-model="newUserName"
            placeholder="New username"
            maxlength="16"
            minlength="3"
          />
          <button
            class="primary-btn"
            :disabled="!canUpdateUsername"
            @click="updateUsername"
          >
            Update Username
          </button>
      </div>

      <ErrorMsg v-if="errormsg" :msg="errormsg" />
      <SuccessMsg v-if="successmsg" :msg="successmsg" />

    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from "../components/ErrorMsg.vue";

import SuccessMsg from "../components/SuccessMsg.vue";
export default {
  name: "ProfileView",
  components: {
    ErrorMsg,
    SuccessMsg,
  },
  data() {
    return {
      userName: "", 
      userPhoto: null, 
      newUserName: "", 
      newPhoto: null, 
      errormsg: null,
      successmsg: null, 
    };
  },
  computed: {
    canUpdateUsername() {
      return (
        this.newUserName &&
        this.newUserName !== this.userName &&
        this.newUserName.length >= 3 &&
        this.newUserName.length <= 16
      );
    },
  },

  methods: {
    async fetchUserProfile() {
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          this.$router.push({ path: "/" });
          return;
        }
        const res = await axios.get("/users/photo", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.userName = localStorage.getItem("name");
        this.userPhoto = res.data.photo
          ? `data:image/jpeg;base64,${res.data.photo}`
          : null;
      } catch (e) {
        console.error("Failed to fetch user profile:", e);
        this.errormsg = "Failed to load user profile. Please try again later.";
      }
    },
    handlePhotoUpload(e) {
      const file = e.target.files[0];
      if (!file) return;
      this.newPhoto = file;
    },
    async updatePhoto() {
      if (!this.newPhoto) return;
      this.errormsg = null;
      try {
        const token = localStorage.getItem("token");
        const form = new FormData();
        form.append("photo", this.newPhoto);
        await axios.put("/users/photo", form, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.successmsg = "Username updated successfully!";
        setTimeout(() => { this.successmsg = null; }, 3000);
        this.newPhoto = null;
        this.fetchUserProfile(); 
      } catch (error) {
        console.error("Failed to update photo:", error);
        this.errormsg = "Failed to update photo. Please try again.";
      }
    },
    async updateUsername() {
      try {
        const token = localStorage.getItem("token");
        const response = await axios.put(
          "/users/name",
          { name: this.newUserName },
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        this.successmsg = "Username updated successfully!";
        this.errormsg = null; // Clear any previous error
        setTimeout(() => { this.successmsg = null; }, 3000);
        localStorage.setItem("name", this.newUserName);
        this.userName = response.data.name;
        this.newUserName = response.data.name;
      } catch (error) {
        console.error("Failed to update username:", error);
        this.errormsg = "Failed to update username. Please try again.";
      }
    },
    
    logOut() {
      localStorage.clear();
      this.$router.push({ path: "/" });
    },

  },
  mounted() {
    this.fetchUserProfile();
  },
};
</script>

<style scoped>
.profile-page {
  min-height: calc(100vh - 60px);
  justify-content: center;
  padding-top: 40px;
  background: #f4f7fb;
}

.profile-topbar-custom {
  background: #f4f7fb;
  color: black;
  padding: 25px 20px; 
    
}

.username-header {
  font-size: 1.5rem;
  font-weight: 500;
  margin: 0;
  color: #2c3e50;
}

.btn-chat-tertiary {
  background-color: transparent;
  color: #7f8c8d;
  border: 1px solid #bdc3c7;
    
  padding: 10px 18px; 
  font-weight: 500;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s, color 0.2s, border-color 0.2s;
}

.btn-chat-tertiary:hover {
  background-color: #ecf0f1;
  color: #34495e;
  border-color: #34495e;
}

.btn-chat-tertiary:active {
  transform: scale(0.97);
}

/* Card */
.profile-card {
  width: 100%;
  max-width: 420px;
  margin: 0 auto;
  background: #ffffff;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.avatar-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.avatar {
  width: 160px;
  height: 160px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid #4a90e2;
}

.avatar.placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f0f0;
  color: #888;
  font-size: 14px;
}


.form-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 18px;
}

.form-group input[type="text"],
.form-group input[type="file"] {
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid #ccc;
  font-size: 14px;
}

.primary-btn {
  padding: 12px;
  font-size: 15px;
  font-weight: 600;
  border-radius: 10px;
  border: none;
  background-color: #4a90e2;
  color: white;
  cursor: pointer;
  transition: background-color 0.2s ease, transform 0.1s ease;
}

.primary-btn:disabled {
  background: #aac7ed;
    cursor: not-allowed;
}
</style>