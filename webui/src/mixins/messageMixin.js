import axios from 'axios';

export const messageMixin = {
    data() {
      return {
        newMessage: '',
        messages: []
      };
    },
    methods: {
      // Send message
      async sendMessage() {
        if (!this.newMessage.trim()) return;

        const token = localStorage.getItem('access_token')?.trim();
        if (!token) {
          console.error('Token missing! User not authorized.');
          return;
        }

        const messageData = {
          chat_id: Number(this.chatId),
          content: this.newMessage.trim()
        };

        try {
          const response = await axios.post("http://localhost:8080/api/messages", messageData, {
            headers: { Authorization: `Bearer ${token}` }
          });

          this.messages.push({
            ID: response.data.ID,
            Content: this.newMessage,
            SenderID: response.data.SenderID,
            isMine: true,
          });

          this.newMessage = ''; // Clear input field
        } catch (error) {
          console.error('Error sending message:', error);
        }
      },

      // Delete message
      async deleteMessage(messageId) {
        const token = localStorage.getItem('access_token')?.trim();
        if (!token) {
          console.error('Token missing! User not authorized.');
          return;
        }

        try {
          await axios.delete(`http://localhost:8080/api/messages/${messageId}`, {
            headers: { Authorization: `Bearer ${token}` }
          });

          this.messages = this.messages.filter(message => message.ID !== messageId);
        } catch (error) {
          console.error('Error deleting message:', error);
        }
      },

      // Add comment
      async commentMessage(messageId, commentText) {
        if (!commentText?.trim()) return;

        const token = localStorage.getItem('access_token')?.trim();
        if (!token) {
          console.error('Token missing! User not authorized.');
          return;
        }

        try {
          await axios.put(`http://localhost:8080/api/messages/${messageId}/comment`, {
            comment: commentText
          }, {
            headers: { Authorization: `Bearer ${token}` }
          });

          const message = this.messages.find(msg => msg.ID === messageId);
          if (message) {
            message.Comment = commentText;
          }
        } catch (error) {
          console.error('Error adding comment:', error);
        }
      },

      // Remove comment
      async uncommentMessage(messageId) {
        const token = localStorage.getItem('access_token')?.trim();
        if (!token) {
          console.error('Token missing! User not authorized.');
          return;
        }

        try {
          await axios.delete(`http://localhost:8080/api/messages/${messageId}/comment`, {
            headers: { Authorization: `Bearer ${token}` }
          });

          const message = this.messages.find(msg => msg.ID === messageId);
          if (message) {
            message.Comment = null;
          }
        } catch (error) {
          console.error('Error deleting comment:', error);
        }
      },

      // Forward message
      async forwardMessage(message) {
        if (!message.forwardChatName.trim()) {
          alert('Enter the chat name!');
          return;
        }

        const token = localStorage.getItem('access_token')?.trim();
        if (!token) {
          console.error('Token missing! User not authorized.');
          return;
        }

        try {
          await axios.post(
            "http://localhost:8080/api/messages/forward",
            {
              message_id: message.ID,
              chat_name: message.forwardChatName
            },
            {
              headers: { Authorization: `Bearer ${token}` }
            }
          );

          alert('Message forwarded!');
          message.forwardChatName = ''; // Clear field
        } catch (error) {
          alert("Chat doesnt exist!");
          console.error('Error forwarding message:', error);
        }
      }
    }
};
