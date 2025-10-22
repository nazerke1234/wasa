<template>
    <div class="new-group-form">
      <h1 class="main-heading">Create New Group</h1>
      <div class="input-container">
        <label for="group-title">Group Title:</label>
        <input
          id="group-title"
          v-model="groupTitle"
          class="text-input"
          type="text"
          placeholder="Enter group title"
        />
      </div>
      <form @submit.prevent="findUsers" class="search-container">
        <input
          id="user-search"
          v-model="searchTerm"
          class="search-input"
          type="text"
          placeholder="Find users by username"
        />
        <button class="search-btn" type="submit">Find</button>
      </form>
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      <div v-if="isLoading">
        <LoadingIndicator />
      </div>
      <div v-if="!isLoading && displayResults" class="results-container">
        <h2 class="results-heading">Found Users:</h2>
        <template v-if="foundUsers.length > 0">
          <div v-for="user in foundUsers" :key="user.id" class="user-item">
            <h5 class="user-username">@{{ user.name }}</h5>
            <button
              class="select-btn"
              @click="includeUser(user)"
              :disabled="isUserIncluded(user)"
            >
              {{ isUserIncluded(user) ? 'Selected' : 'Select' }}
            </button>
          </div>
        </template>
        <template v-else>
          <p class="empty-results">No users found for "{{ previousSearch }}"</p>
        </template>
      </div>
      <div v-if="includedUsers.length > 0" class="members-container">
        <h2 class="members-heading">Group Members:</h2>
        <div class="members-list">
          <span v-for="user in includedUsers" :key="user.id" class="member-tag">
            @{{ user.name }}
            <button @click="excludeUser(user)" class="deselect-btn">Remove</button>
          </span>
        </div>
      </div>
      <div class="input-container">
        <label for="group-picture">Group Picture:</label>
        <input
          id="group-picture"
          ref="imageInput"
          type="file"
          @change="processImageSelection"
          accept="image/*"
        />
        <img v-if="imagePreview" :src="imagePreview" class="preview-img" alt="Group picture preview"/>
      </div>
      <button class="submit-btn" @click="submitGroup" :disabled="!isFormValid">
        Create Group
      </button>
    </div>
  </template>

  <script>
  import axios from "../services/axios";
  import LoadingIndicator from "../components/LoadingSpinner.vue";
  
  export default {
    name: "CreateGroupView",
    components: {
      LoadingIndicator,
    },
    data() {
      const authToken = localStorage.getItem("token");
      if (!authToken) {
        this.$router.push({ path: "/" });
        return;
      }
      return {
        groupTitle: "",
        searchTerm: "",
        previousSearch: "",
        foundUsers: [],
        isLoading: false,
        displayResults: false,
        errorMessage: "",
        includedUsers: [],
        imagePreview: null,
        selectedImage: null,
      };
    },
    computed: {
      isFormValid() {
        return (
          this.groupTitle.trim() !== "" &&
          this.includedUsers.length > 0 &&
          this.selectedImage !== null
        );
      },
    },
    methods: {
      async findUsers() {
        if (!this.searchTerm.trim()) {
          this.errorMessage = "Please provide a valid search term.";
          this.displayResults = false;
          return;
        }
        this.isLoading = true;
        this.errorMessage = "";
        this.foundUsers = [];
        this.displayResults = false;
        try {
          const response = await axios.get(`/search`, {
            params: { username: this.searchTerm },
          });
          this.foundUsers = response.data.filter(user => user.id !== localStorage.getItem("token"));
          this.previousSearch = this.searchTerm;
          this.displayResults = true;
        } catch (err) {
          const errorCode = err.response?.status;
          const errorDetails = err.response?.data?.message || "Unable to retrieve users.";
          this.errorMessage = `Error ${errorCode}: ${errorDetails}`;
        } finally {
          this.isLoading = false;
        }
      },
      includeUser(user) {
        if (!this.isUserIncluded(user)) {
          this.includedUsers.push(user);
        }
      },
      isUserIncluded(user) {
        return this.includedUsers.some((u) => u.id === user.id);
      },
      excludeUser(user) {
        this.includedUsers = this.includedUsers.filter((u) => u.id !== user.id);
      },
      processImageSelection(event) {
        const imageFile = event.target.files[0];
        if (imageFile) {
          this.selectedImage = imageFile;
          const fileReader = new FileReader();
          fileReader.onload = (e) => {
            this.imagePreview = e.target.result;
          };
          fileReader.readAsDataURL(imageFile);
        } else {
          this.imagePreview = null;
          this.selectedImage = null;
        }
      },
      async submitGroup() {
        if (!this.isFormValid) {
          alert("Please complete all required fields and choose a group image.");
          return;
        }
        this.isLoading = true;
        this.errorMessage = "";
        const formPayload = new FormData();
        formPayload.append("name", this.groupTitle);
        formPayload.append("image", this.selectedImage);
        formPayload.append("members", JSON.stringify([...this.includedUsers.map(u => u.id), localStorage.getItem("token")]));
        try {
          await axios.post(`/groups`, formPayload, {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          });
          alert("Group successfully created!");
          this.$router.push(`/home`);
        } catch (err) {
          const errorCode = err.response?.status;
          const errorDetails = err.response?.data?.message || "Group creation failed.";
          this.errorMessage = `Error ${errorCode}: ${errorDetails}`;
        } finally {
          this.isLoading = false;
        }
      },
    },
  };
  </script>
  <style scoped>
  .new-group-form {
    text-align: center;
    padding: 20px;
    max-width: 600px;
    margin: 0 auto;
  }
  
  .main-heading {
    font-size: 28px;
    font-weight: bold;
    margin-bottom: 20px;
    color: #333;
  }
  
  .input-container {
    margin-bottom: 20px;
  }
  
  .text-input, .search-input {
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
    font-size: 16px;
    width: 100%;
    box-sizing: border-box;
  }
  
  .search-container {
    display: flex;
    justify-content: center;
    margin-bottom: 20px;
  }
  
  .search-btn, .submit-btn {
    padding: 10px 20px;
    background-color: #007bff;
    color: #fff;
    border: none;
    border-radius: 5px;
    font-size: 16px;
    cursor: pointer;
    margin-left: 10px;
  }
  
  .search-btn:hover, .submit-btn:hover {
    background-color: #0056b3;
  }
  
  .error-message {
    background-color: #f8d7da;
    color: #842029;
    border: 1px solid #f5c2c7;
    border-radius: 5px;
    padding: 10px;
    margin: 20px 0;
    text-align: center;
  }
  
  .results-container {
    margin-top: 20px;
    max-height: 300px; 
    overflow-y: auto;
    border: 1px solid #ccc;
    border-radius: 5px;
    padding: 10px;
  }
  
  .results-heading {
    font-size: 24px;
    font-weight: bold;
    color: #444;
    margin-bottom: 10px;
  }
  
  .user-item {
    padding: 10px;
    margin: 10px 0;
    background-color: #f9f9f9;
    border: 1px solid #ccc;
    border-radius: 5px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .user-username {
    font-size: 18px;
    color: #007bff;
  }
  
  .select-btn, .deselect-btn {
    padding: 5px 10px;
    background-color: #28a745;
    color: #fff;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 14px;
  }
  
  .select-btn:hover, .deselect-btn:hover {
    background-color: #218838;
  }
  
  .select-btn[disabled], .deselect-btn[disabled] {
    background-color: #ccc;
    cursor: not-allowed;
  }
  
  .members-container {
    margin-top: 20px;
  }
  
  .members-heading {
    font-size: 20px;
    font-weight: bold;
    color: #444;
    margin-bottom: 10px;
  }
  
  .members-list {
    display: flex;
    flex-wrap: wrap;
  }
  
  .member-tag {
    padding: 5px 10px;
    background-color: #e9ecef;
    border-radius: 5px;
    margin: 5px;
    font-size: 16px;
  }
  
  .preview-img {
    max-width: 200px;
    max-height: 200px;
    margin-top: 10px;
  }
  @media (max-width: 768px) {
  .results-container {
    max-height: 200px;
  }
}
  </style>

