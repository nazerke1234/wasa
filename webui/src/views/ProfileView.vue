<template>
  <div>
    <div class="profile-header-section">
      <h1 class="profile-title">{{ userDisplayName }}, your profile details</h1>
      <div class="action-controls">
        <div class="control-group">
          <button type="button" class="control-btn secondary" @click="reloadProfile">Reload</button>
          <button type="button" class="control-btn secondary" @click="signOut">Sign Out</button>
        </div>
        <div class="control-group">
          <button type="button" class="control-btn primary" @click="createGroup">Create Group</button>
        </div>
      </div>
    </div>
    
    <div class="profile-content">
      <div class="profile-card">
        <div class="avatar-section">
          <img v-if="userAvatar" :src="userAvatar" alt="User Avatar" class="avatar-image" />
          <div v-else class="avatar-placeholder">No Avatar</div>
        </div>
        <div class="profile-details">
          <h1 class="display-name">{{ userDisplayName }}</h1>
          <div class="update-section">
            <input
              v-model="updatedUsername"
              placeholder="Enter new display name"
              maxlength="16"
              minlength="3"
              class="name-input"
            />
            <button
              class="update-btn"
              @click="modifyUsername"
              :disabled="!updatedUsername || updatedUsername === userDisplayName"
            >
              Change Display Name
            </button>
          </div>
          <div class="avatar-update-section">
            <input type="file" @change="processAvatarSelection" accept="image/*" class="file-input" />
            <button class="update-btn" @click="modifyAvatar" :disabled="!newAvatarFile">
              Update Avatar
            </button>
          </div>
        </div>
      </div>
      <ErrorMsg v-if="errorMessage" :msg="errorMessage" />
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "ProfileView",
  components: {
    ErrorMsg,
  },
  data() {
    return {
      userDisplayName: "", 
      userAvatar: null, 
      updatedUsername: "", 
      newAvatarFile: null, 
      errorMessage: null, 
    };
  },
  methods: {
    async loadUserProfile() {
      try {
        const authToken = localStorage.getItem("token");
        if (!authToken) {
          this.$router.push({ path: "/" });
          return;
        }
        const response = await axios.get("/users/photo", {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });
        const { photo } = response.data;
        this.userDisplayName = localStorage.getItem("name");
        this.userAvatar = photo ? `data:image/jpeg;base64,${photo}` : null;
      } catch (error) {
        console.error("Unable to load user profile:", error);
        this.errorMessage = "Unable to load profile information. Please try again later.";
      }
    },
    processAvatarSelection(event) {
      const selectedFile = event.target.files[0];
      if (selectedFile) {
        this.newAvatarFile = selectedFile;
      }
    },
    async modifyAvatar() {
      if (!this.newAvatarFile) return;
      try {
        const authToken = localStorage.getItem("token");
        const imageData = new FormData();
        imageData.append("photo", this.newAvatarFile);
        await axios.put("/users/photo", imageData, {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });
        alert("Avatar updated successfully!");
        this.newAvatarFile = null;
        this.loadUserProfile();
      } catch (error) {
        console.error("Unable to update avatar:", error);
        this.errorMessage = "Unable to update avatar. Please try again.";
      }
    },
    async modifyUsername() {
      if (!this.updatedUsername || this.updatedUsername === this.userDisplayName) return;
      try {
        const authToken = localStorage.getItem("token");
        const response = await axios.put(
          "/users/name",
          { name: this.updatedUsername },
          {
            headers: {
              Authorization: `Bearer ${authToken}`,
            },
          }
        );
        alert("Display name updated successfully!");
        localStorage.setItem("name", this.updatedUsername);
        this.userDisplayName = response.data.name;
        this.updatedUsername = response.data.name;
        this.errorMessage = null;
      } catch (error) {
        console.error("Unable to update display name:", error);
        if (error.response && error.response.status === 409) {
          this.errorMessage = "This username is already taken. Please choose a different one.";
        } else if (error.response && error.response.status === 400) {
          this.errorMessage = "Invalid username. Must be 3-16 characters and contain only letters, numbers, and underscores.";
        } else {
          this.errorMessage = "Unable to update display name. Please try again.";
        }
      }
    },
    reloadProfile() {
      this.loadUserProfile();
    },
    signOut() {
      localStorage.clear();
      this.$router.push({ path: "/" });
    },
    createGroup() {
      this.$router.push({ path: "/new-group" });
    }
  },
  mounted() {
    this.loadUserProfile();
  },
};
</script>

<style scoped>
.profile-header-section {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  align-items: center;
  padding-top: 1rem;
  padding-bottom: 0.5rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid #e0e0e0;
}

.profile-title {
  font-size: 1.75rem;
  font-weight: 600;
  color: #333;
}

.action-controls {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.control-group {
  display: flex;
  gap: 0.5rem;
}

.control-btn {
  padding: 0.375rem 0.75rem;
  border: 1px solid;
  border-radius: 0.375rem;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.control-btn.secondary {
  background-color: transparent;
  color: #6b7280;
  border-color: #6b7280;
}

.control-btn.secondary:hover {
  background-color: #6b7280;
  color: white;
}

.control-btn.primary {
  background-color: transparent;
  color: #3b82f6;
  border-color: #3b82f6;
}

.control-btn.primary:hover {
  background-color: #3b82f6;
  color: white;
}

.profile-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1.5rem;
}

.profile-card {
  display: flex;
  align-items: flex-start;
  gap: 2rem;
  width: 100%;
  max-width: 800px;
  background: #ffffff;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border: 1px solid #e9ecef;
}

.avatar-section {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid #e9ecef;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  flex-shrink: 0;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  color: #9ca3af;
  font-size: 0.875rem;
  font-weight: 500;
}

.profile-details {
  flex: 1;
  min-width: 0;
}

.display-name {
  margin: 0 0 1.5rem 0;
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
  background: linear-gradient(135deg, #3b82f6 0%, #1e40af 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.update-section,
.avatar-update-section {
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.name-input,
.file-input {
  padding: 0.75rem 1rem;
  border: 2px solid #e5e7eb;
  border-radius: 0.5rem;
  flex: 1;
  max-width: 300px;
  font-size: 0.875rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.name-input:focus,
.file-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.update-btn {
  padding: 0.75rem 1.5rem;
  background-color: transparent;
  border: 2px solid #3b82f6;
  color: #3b82f6;
  border-radius: 0.5rem;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.875rem;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.update-btn:hover:not(:disabled) {
  background-color: #3b82f6;
  color: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(59, 130, 246, 0.3);
}

.update-btn:disabled {
  border-color: #d1d5db;
  color: #d1d5db;
  cursor: not-allowed;
  background-color: transparent;
  transform: none;
  box-shadow: none;
}

@media (max-width: 768px) {
  .profile-header-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .action-controls {
    width: 100%;
    justify-content: space-between;
  }
  
  .profile-card {
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 1.5rem;
  }
  
  .update-section,
  .avatar-update-section {
    flex-direction: column;
    align-items: stretch;
  }
  
  .name-input,
  .file-input {
    max-width: none;
  }
}
</style>


