<template>
    <div class="group-create-container">
        <h1 class="page-title">Create Group</h1>
        
        <div v-if="error" class="error-box"> 
            <ErrorMsg :msg="error" />
        </div>
        <LoadingSpinner v-if="loading" />

        <div class="form-section">
            <label for="group-name" class="input-label">Group Name (Required):</label>
            <input
                id="group-name"
                v-model="groupName"
                class="app-input-field"
                type="text"
                placeholder="Enter a title for your group (e.g., Team Project)"
            />
        </div>

        <div class="form-section search-section">
            <h2 class="section-title">Add Members</h2>
            <form @submit.prevent="searchUsers" class="search-form">
                <input
                    id="username"
                    v-model="query"
                    class="app-input-field search-box"
                    type="text"
                    placeholder="Search by username to add members"
                />
                <button 
                    class="btn-secondary search-button" 
                    type="submit" 
                    :disabled="loading"
                >
                    Search
                </button>
            </form>
        </div>

        <div v-if="!loading && showResults" class="results-section app-card">
            <h3 class="results-title">Results for "{{ lastQuery }}"</h3>
            <template v-if="users.length > 0">
                <div v-for="user in users" :key="user.id" class="user-card">
                    <span class="user-name">@{{ user.name }}</span>
                    <button
                        class="btn-small-tertiary"
                        @click="addUserToGroup(user)"
                        :disabled="isUserAdded(user)"
                    >
                        {{ isUserAdded(user) ? 'Added' : 'Add' }}
                    </button>
                </div>
            </template>
            <template v-else>
                <p class="no-results">No users found matching "{{ lastQuery }}"</p>
            </template>
        </div>

        <div v-if="selectedUsers.length > 0" class="selected-users-section">
            <h2 class="section-title selected-title">Selected Members ({{ selectedUsers.length }}):</h2>
            <div class="selected-users-tags">
                <span v-for="user in selectedUsers" :key="user.id" class="selected-user-tag">
                    @{{ user.name }}
                    <button @click="removeUserFromGroup(user)" class="remove-tag-button">
                        &times;
                    </button>
                </span>
            </div>
        </div>

        <div class="form-section image-upload-section">
            <label for="group-image" class="input-label">Group Image:</label>
            <input
                id="group-image"
                ref="fileInput"
                type="file"
                @change="handleImageUpload"
                accept="image/jpeg, image/png"
            />
            <div class="image-preview-wrapper" v-if="previewImage">
                 <img :src="previewImage" class="preview-image" alt="Group Image Preview"/>
            </div>
           
        </div>

        <div class="action-footer">
            <button 
                class="btn-primary create-button" 
                @click="createGroup" 
                :disabled="!canCreateGroup || loading"
            >
                {{ loading ? 'Creating...' : 'Create Group' }}
            </button>
        </div>
    </div>
</template>

<script>

import axios from "../services/axios";


export default {
    name: "CreateGroupView",
    
    data() {
      const token = localStorage.getItem("token");
      if (!token) {
        this.$router.push({ path: "/" });
      }
      return {
        groupName: "",
        query: "",
        lastQuery: "",
        users: [],
        loading: false,
        showResults: false,
        error: "",
        selectedUsers: [],
        previewImage: null,
        file: null,
        token: token, 
      };
    },
    computed: {

      isNameValid() {
          return this.groupName.trim().length >= 3;
      },
      
      canCreateGroup() {
        return (
          this.isNameValid &&
          this.selectedUsers.length > 0 
        );
      },
    },
    methods: {
      async searchUsers() {
        if (!this.query.trim()) {
          this.error = "Please enter a valid search query.";
          this.showResults = false;
          return;
        }
        this.loading = true;
        this.error = "";
        this.users = [];
        this.showResults = false;
        try {
          const response = await axios.get(`/search`, {
            params: { username: this.query },
            headers: { Authorization: `Bearer ${this.token}` }, 
          });
          
          
          const myId = localStorage.getItem("token"); 
          this.users = response.data.filter(user => user.id !== myId); 
          
          this.lastQuery = this.query;
          this.showResults = true;
        } catch (err) {
          const status = err.response?.status;
          const reason = err.response?.data?.message || "Failed to fetch users.";
          this.error = `Status ${status}: ${reason}`;
        } finally {
          this.loading = false;
        }
      },
      addUserToGroup(user) {
        if (!this.isUserAdded(user)) {
          this.selectedUsers.push(user);
        }
      },
      isUserAdded(user) {
        return this.selectedUsers.some((u) => u.id === user.id);
      },
      removeUserFromGroup(user) {
        this.selectedUsers = this.selectedUsers.filter((u) => u.id !== user.id);
      },
      handleImageUpload(event) {
        const file = event.target.files[0];
        if (file) {
          this.file = file;
          const reader = new FileReader();
          reader.onload = (e) => {
            this.previewImage = e.target.result;
          };
          reader.readAsDataURL(file);
        } else {
          this.previewImage = null;
          this.file = null;
        }
      },
      async createGroup() {
        if (!this.canCreateGroup) {
          this.error = "Please provide a valid Group Name (min 3 chars), select members, and upload a Group Image.";
          return;
        }
        this.loading = true;
        this.error = "";

        
        const memberIds = [
            ...this.selectedUsers.map(u => u.id), 
            localStorage.getItem("token") 
        ];

        const formData = new FormData();
        formData.append("name", this.groupName);
        formData.append("image", this.file);
        
        formData.append("members", JSON.stringify(memberIds)); 
        
        try {
          await axios.post(`/groups`, formData, {
            headers: {
              'Content-Type': 'multipart/form-data', 
              'Authorization': `Bearer ${this.token}` 
            }
          });
          
          alert("Group created successfully! Redirecting to Home.");
          this.$router.push(`/home`);
        } catch (err) {
          const status = err.response?.status;
          const reason = err.response?.data?.message || "Failed to create group.";
          this.error = `Error ${status}: ${reason}. Please check group name uniqueness and member list.`;
        } finally {
          this.loading = false;
        }
      },
    },
  };
</script>

<style scoped>


.group-create-container {
  max-width: 700px;
  margin: 30px auto;
  padding: 30px;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 25px;
  color: #4a90e2;
  text-align: center;
}

.form-section {
    margin-bottom: 25px;
}

.section-title {
    font-size: 1.4rem;
    font-weight: 600;
    color:  #34495e;
    margin-bottom: 15px;
}

.input-label {
    display: block;
    text-align: left;
    margin-bottom: 8px;
    font-weight: 500;
    color:  #34495e;
}

/* Reusing standard input field style */
.app-input-field {
  padding: 12px;
  border: 1px solid  #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  width: 100%;
  box-sizing: border-box;
  transition: border-color 0.2s;
}
.app-input-field:focus {
  border-color: #4a90e2;
  outline: none;
  box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.2);
}


.search-form {
  display: flex;
  gap: 10px;
}

/* Standard Button Styles */
.btn-primary, .btn-secondary {
    padding: 12px 20px;
    font-weight: 600;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: background-color 0.2s, opacity 0.2s;
    white-space: nowrap;
}

/* Search button uses secondary color (Green) */
.btn-secondary {
    background-color: #28a745; 
    color: white;
}
.btn-secondary:hover:not(:disabled) {
    background-color: #1f8b36;
}

.results-section {
  padding: 15px;
  border: 1px solid  #e0e0e0;
  border-radius: 8px;
  max-height: 250px;
  overflow-y: auto;
  margin-top: 15px;
  background-color: #fcfcfc;
}

.user-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #f0f0f0;
}
.user-card:last-child {
    border-bottom: none;
}

.user-name {
    font-weight: 500;
    color:  #34495e;
}

/* Tertiary small button for adding */
.btn-small-tertiary {
    padding: 6px 12px;
    font-size: 0.85rem;
    font-weight: 500;
    background-color: #e6f1fc; 
    color: #4a90e2;
    border: 1px solid #4a90e2;
    border-radius: 6px;
    cursor: pointer;
}
.btn-small-tertiary:hover:not(:disabled) {
    background-color: #c0dcfc;
}
.btn-small-tertiary:disabled {
    background-color: #ccc;
    color: #666;
    border-color: #aaa;
    cursor: not-allowed;
}
.no-results {
    color: #95a5a6;
    font-style: italic;
    text-align: center;
    padding: 10px 0;
}


.selected-users-section {
    margin-top: 30px;
    padding-top: 15px;
    border-top: 1px dashed  #e0e0e0;
}

.selected-users-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
}

.selected-user-tag {
    display: inline-flex;
    align-items: center;
    padding: 5px 10px;
    background-color: #e9ecef;
    border-radius: 20px;
    font-size: 0.9rem;
    color: #34495e;
}

.remove-tag-button {
    background: none;
    border: none;
    color: #888;
    margin-left: 5px;
    padding: 0;
    line-height: 1;
    font-size: 1.2rem;
    cursor: pointer;
    transition: color 0.2s;
}
.remove-tag-button:hover {
    color: #4a90e2;
}



.image-upload-section input[type="file"] {
    border: 1px solid  #e0e0e0;
    border-radius: 8px;
    padding: 10px;
    width: 100%;
    box-sizing: border-box;
    background-color: #f7f9fb;
}

.image-preview-wrapper {
    margin-top: 15px;
    border: 2px dashed  #e0e0e0;
    width: 150px;
    height: 150px;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: 8px;
    overflow: hidden;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.action-footer {
    margin-top: 40px;
    padding-top: 20px;
    border-top: 1px solid  #e0e0e0;
}

.create-button {
  width: 100%;
  background-color: #4a90e2;
  color: white;
}

.btn-primary:disabled {
    background-color: #ccc;
    cursor: not-allowed;
    opacity: 0.8;
}

.error-box {
    background-color: #f8d7da;
    color: #842029;
    border: 1px solid #f5c2c7;
    border-radius: 8px;
    padding: 10px;
    margin: 20px 0;
    text-align: center;
}
</style>