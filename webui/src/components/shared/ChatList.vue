<template>
  <div class="main-container">
    <router-link to="/search-users" class="button search-users">Search Users</router-link>
    <router-link to="/profile" class="button profile">Profile</router-link>

    <div class="chat-container">
    <h2>Your Chats</h2>
    <ul class="chat-list">
      <li v-for="chat in chats" :key="chat.ID" class="chat-item">
        <!-- Логика для изменения ссылки в зависимости от типа чата -->
        <a :href="getChatLink(chat)" class="chat-link">
          <!-- Отображаем фото чата или пользователя -->
          <img v-if="chat.PhotoUrl" :src="chat.PhotoUrl" alt="Chat photo" width="50" height="50" class="chat-img"/>
          <img v-else src="data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBxISEhUTEBIVFhUXFxUVFRYVFhUQFxgVFRIWGBcWFhUYHSggGB0lGxcWITEhJSkrLi4uFx8zODMtNygtLisBCgoKDQ0NDw0NFSsdHxkrKysrKystLSsrKysrLSsrKysrKystKysrKysrKysrKysrKysrKysrKysrKysrKysrK//AABEIAOEA4QMBIgACEQEDEQH/xAAbAAEAAgMBAQAAAAAAAAAAAAAAAQYCBAUHA//EAD4QAAECAgUICAUDBAIDAAAAAAEAEQIhAzEyQVEEBRJCYaHB8AYTIlJicYGRcoKx0eEjM6JDU5KyY/EUc8L/xAAVAQEBAAAAAAAAAAAAAAAAAAAAAf/EABQRAQAAAAAAAAAAAAAAAAAAAAD/2gAMAwEAAhEDEQA/APZQAAwnCazghAZjKEVHHmaAhnEobxihIZzZNQw5mgkkkuZRCyMUBILicRrGHMkILsZxGycEALsJRXnHmSCAAAQJwm0cEIDMbIqixPLoCCHEoRaGKEhnNm4YHl/dBJJLEyiFkYoCXcWrxgOWQgiRnEbJwQAuwtXnn0QQAzgTBtHBCAzGxdFt5dSJzEgLQxXxyjKqOAPSRAQ3Qkz9r70H2JdiZEWRijl3Fu+HZyy4mUdJKMPoiKM3GwB7z3Ln0vSSmNkQwnFtI75bkFqEnAmDaOCMGY2LotqpUWeKcv8AqEPWwhh+gXxOX0p/qR/5FBfDNtKRFnajl3174diof/n0v9yP/KL7r6wZ3pwXFLE+1ovqEF2En0Zg2tnM0YM2p3tqqdD0iphXoxC+Wi/my6GT9JYCWpIDDDhCRGHxuKDumbPJrO3mSOXfX7uxa+TZdR0lmMRHVAM4fMGeHsthi7a/e2IAk7Te1s5moYMz9jvbVIvaTWtvM1Dhn1O7tQSZs8ms7eZe6OXfX7uxDJtKb2dnMvZGLtr97YgCTtN7WzmahgzP2O9tUibtJrW3mahwz6nd2oBmzyazt5kpcu+v3diGTPN7OzmSMXZ+33tiDLraTuonV0nfHPoiDEkkuZRCoYoCXcTiNcOHMkIILGcRqOCAF2Fq848yQAAAwnCazgoIDMZQ3RYqQQziUItDFQSGc2bhhzNBJM3MohUMUBLuJxXw4cy90ILsZxGycEm7AtEKzc3LIIAZwJg2jgtbLcto6INSRMK4WnFF5BcrOufwHgoPKKKsPsx81XKSMxF4iSTWTMoOvnDpDSRn9MaAF4nF73Ljxxkl4iScSXPuVCICIiAiIgIiICIiACupkefaWAaMXbhwMj6Rfd1y0QXfIM50dNUWiFUJkT9/Rbrl31u7sXnkJImJEVESXezZ0gI7NMdgpLx8Qv8ANBZRJ2m9rw8z9lDBmfs97akEQiDwGVcRrEQ2b0cM+phtQSZs8ms+LmSOXfW7uxDdpTezs5kjF218diAJO03tbOZ+yhgzP2e9tUi9pNa28zUOGfUw2oI6qDvonWUfdPPqiCQAAwLwmuLBCAzEtCKosUDM8ModYYoSGc2dUXjmaCSXmZRCoYoCXcTiNcOCEF2No2Tgk3YFohWbm5ZBiSADPs1xRHVVVz1nk0n6dGWoxfUYvPAbEz7nbrD1dFKjFfiIv8lx0BERAREQEREBERAREQEREBERAREQdHNWdYqEgFzA76OG2H7K30FOIwKSAgkiQGC8/XQzPnOKhi8JrvbaEF0Enab2vDzNQwZn7Pe2qKKMRB4DJniN0Q2b1Lhn1ML3QSZ1yaz4uZI5d27Xd2IbtKb2NnnuRi7a+Nzc7EE9bH3EU9XSd4c+iIMSXLkNEKocUBm4DxGuHBCC7RWtU4IHdha1jjzJBADSBcGs4Lg9JM5MOpgPxRDDBdXOGVijo4oxVVonWiNQ5wVHpIzESSXJLk7SgxREQEREBERAREQEREBERAREQEREBERAREQdzo5nFiKKMyJeDZF3fI/XzVoebt2u7xXnauuZMv66Af3BKI8fUIN8Sdpva8PPBQwZn7Pe4KRfo3W9vlvUOGfUwvdBHVwd9E06Punn1RUSABIFwa4sEa4lgKosUBDPDZ1gvlldMIKOKOKcIBMI23KCtdJ8uMdIIKhBX8Rr9qvdcZTHGSSTWSSfMlyoQEREBERAREQEQBWLNnR6qKn9IB/9HgEHAoqKKItDCYjsBK3YMyZQf6becUI4q40VFDCGhAAwAZZIKbHmPKB/TfyihPFaVPk8cFuEw+YIV/URwghiARgZhB56is+cuj8MTxUPZi7uqfLD6KtUlGYSYYgQRIg3IMUREBERAREQF0Mx5X1dKHLQxdmL1Mj6H6lc9EHohnXJrPi54o83btd3itLNGVdZQwRRlzZGyKGU9x9VuzdtfG5kE9bF3EU6FJ3hz6IgxJeZDEVQ4ridKqZqMQvOOKYwEM/qy7Zd2it6puVU6VUj0whNcMIfzM/oyDjIiICIiAiIgIi3M05J1tLDCaq4vhH3kPVB2+jubNECljHaNkYDHzP0XcREBERAREQFy8+ZsFLDpQjtwiXiHdPBdREHniLqdIsk0KVxVH2vV+1vn6rloCIiAiIgIiILD0Syggx0bO4EXlcTvCsTS0X7Pe4KmZgpdGng2vCfIj/pXNwz6mF7oHVQ99QmnR4Hf90QSzSBcGuLBUrPcT09JN2ID+UIHBXUM3Zsawv5qVEziXpaRu/F/sUGuiIgIiICIiArB0Sop0kXww+7k/QKvqy9Ej2KT4h/qg7yIiAiIgIiICIiDi9KqJ6OGLCLdED9gqsrd0mP6B+KH6qooCIiAiIgIiIPrkkejSQRYRQn2iCv73tPu8V52CvRA/z43MgnrYv7aKdGkxHPoiDF3mQxFUOKoecf3aT44v8AYq+F37VvVNypGeoCKekBr0n9wDxQaSIiAiIgIiIC7nRSmakig70LjzhP2J9lw19clpzRxwxw1wl/uPZBfkWFBTCOERQlwQ4WaAiIgIiICIoiiADmQEydgQcHpZTdmCDEmI+QDD6n2VbW3nTK+tpIo7qofhFX39VqICIiAiIgIiIAXoYEtF5d7gqBk0DxwDGKEe8QCv7hvBhe6B1Y/ub0UaVHgd/3RBLNJ3B1sFU+k9E1M/ehE8Wk/wBFbAzdmxrYrh9K6AmCCMWQTCMWiH3G9BWEREBERAREQEREHVzHnXqjox2D/E4+WKtsEQIBBBBmCJgrz1bubs6UlDZLw3wmr0wQXZFyslz/AEMdomA4RTH+Q/C6NHlEEVmOE+UQKD6IsI6eEVxQjziAWhlOfKGCqLSOEM99SDpKs5+zvp/p0Z7OtF3tg2fVamcs80lLKzB3Qa/M3rmoCIiAiIgIiICIiDfzFRaVPBsOl/iH+rK6vrNPucWVZ6KUI0o4zcBDD8RL8B7qzTfx7mQT1p/t7j9kUtSYjcioxd5szauK184ZP1lHHDfECBD3TWD7gLYLv2rerhzWgd5W9bDmpQedkIup0iyUQUulDZjmPiFobwfVctAREQERbOb8jipoxDD6m4DFBGQ5FHSxaMA8zcBiVasjzNRQQkGERkhojEPpgFt5HksNFCIYBLeTidq+yCqZzzFFA8VE8UOFcQ+6469DWnlua6KlnFCx70PZP59UFIUMrBT9GT/TpAdkQbePstSLo/Ti6E+UQ4oOUyldOHMFOdUDziHBbVB0Zj144R8IMR3sg4S6ebMzR0rGLswYms/COKsGR5moaOYh0jjF2vYVBdBBoUuZ6EwCDRZqohafF71V85ZtjoT2pwmqIVHZsOxXdYU9DDHCYYg4NYQefot7O2bjQxNXCbJ4HatFAREQERbmaMjNLSww3WovIct6oLRmPJuroYQzmLtReF6tzey6DXP8/B0F+j8/43qJN4N7oHVD+5v/ACij9PbvRUZM0nd9bBGudm1sVAZuzY1seakLNOxq481qDTzvkXXUZDNEJwjEj71KkxBixkRIjavQy79q1qqtdJc3sethH/sAuN0XrJ0HAREQTBCSQAHJIAG01BXbNOQChgEOsZxHE4eQXE6L5HpRGkIlDKH4jf6D6qzoCIiAiIgIiICIiAiIgIiIPhluSw0sBgivqOBuIVHyihMERgirBY/dX9cDpTkbgUoEx2YvI1H3l6oK2iIgK39H8h6uj7Uoo2Ji7ouh5xXH6PZsNLFpxDsQGe04K2Sbwb3QTXsb+XPFHvb5OLIbtL5PzuSb+PcyB1v/AB7vwiy/V2bkQYu82ZtXFHabO+rhzxQu/at6uCB3lb1sOakBmk7vrd1RHAIgYSzXxGqIYIGY6NnWxQs07F2PNaCm54zaaGJwDoE9l7vCfuucvQMooRGDDSAEkMBcyrFJmOKCmhhrgMQc4CsiL0Qd/NOTdXQwQ3s584pn6t6LbREBERAREQEREBERAREQEREBfPKaERwRQGqIEe4rX0RB57HCQSDWCQfMGa2815vipo9EShFqK4D77F0MszPHSZTGIQ0Mo4orgIq/Mu8lYckyaCjgEMI/TvOsYsSgzoKGGGEQwjREMh4vvvrX0fWb5OPIQ3aXyfnck38e5kCra/8AHngjar/Pw5KC/R+f8b1Em8G90E9V/wAnPuiw/S271KDJmk7vrYI1zs2tjzwUBm7NjWx5qQs07GrjzWgl3mzNq95HvZ/DhzxQu/atauCB3lbvw5qQGaTu+t3UZ5Oza2OznBQGno2dZCzTsXYvy6AJ3Ns+yKS+ta1fyhcy1rxdzUghEE6kQEREBERAREQEREBEUiZa9BCD223KfK61s8t6iTeDe6Ca9jfy54o+s3ycWQ3aXyfnck38e5kCra/8eeCNc/z8HQX6Pz/jeok3g3ugmvY38ueKPrN8nFkN2l8n53JN/HuZA63/AI934ULN6TYiDF3mzNq4o7TZ31cOeKF37VvVwQO8rethzUgM0nd9bBGudm1seeCgM3Zs62KFmnYuxfl0Eu82ZtXvI97O+rhtbmtC79q1qoHeVu/BuWQGaTu+t3Ua528WOznBQGno2db8IWadi7F+XQS7zqa7vKD3mr1RWOeKkvLStaqB3lbvwblkCK6978PNOdnuoF+jZ1vwhZp2LsX5dARSbtKvU/O5JvK3fg3LIIRSL9GrX/G9QWadi7F+XQTz/wBYqLn3XqTdpV6nB9yTfx7mQKtr/wAfPm5G1X+fggv0atf8b1Em8G90E17G/lzxR72l3OLIbtL5PzuSb+PcyBVtf+PPBGuf5+DoL9H5/wAb1Em8G90E17G/lzxR72l3OKG7S+T87km/j3MgVbX/AI88Ea5/n4IL9H5/xvUSbwb3QT1R/ub/AMosWo9u9EH0pv3IecUov3IvL7IiD50NiJKT9uHz+6IgzprcKmj/AHIvL7Iio+dDYj5uUR/tw+f3UooM6W1Apg/cPl9kRB86GzH6pH+2PPiURBnTWoOcFMP7h8uAREGFFVH68VEf7Y8+JREGdNXB6fUKR+6fLgERUYUNVJ68VB/b9eKIoMqbU9OCy/q+nBEQY0X9T14rD+l68VKIJptT04LM/uenBEQY0NcfrxWEP7XrxUoqNZEREf/Z" alt="Default photo" width="50" height="50" />
          <span class="chat-name">
            <span>{{ chat.Name }}</span>
          </span>
        </a>
      </li>
    </ul>
  </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      chats: []  // Храним чаты
    };
  },
  mounted() {
    this.getChats(); // При загрузке компонента вызываем метод для получения чатов
  },
  methods: {
    getChats() {
      axios.get('http://localhost:8080/api/chats', {
        headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` }
      })
      .then(response => {
        console.log('Chats data:', response.data);  // Выводим полученные данные чатов
        this.chats = response.data;   // Сохраняем чаты из ответа в массив
      })
      .catch(error => {
        console.error('Ошибка при получении чатов:', error);
      });
    },

    // Функция для генерации правильного URL для чата
    getChatLink(chat) {
      if (chat.IsGroup) {
        return `/groupChats/${chat.ID}`;  // Ссылка на групповое сообщение
      } else {
        return `/chats/${chat.ID}`;  // Ссылка на личное сообщение
      }
    }
  }
};
</script>

<style scoped>
/* Общие стили */
.main-container {
  font-family: 'Poppins', sans-serif;
  background: linear-gradient(135deg, #1e3c3c, #2a5d5d);
  color: #333;
  margin: 0;
  padding: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  position: relative;
}

/* Блок чатов */
.chat-container {
  width: 400px;
  background: white;
  padding: 20px;
  border-radius: 15px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  text-align: center;
}

/* Список чатов */
.chat-list {
  list-style: none;
  padding: 0;
  margin-top: 20px;
}

.chat-item {
  display: flex;
  align-items: center;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 8px;
  margin: 5px 0;
  transition: background 0.3s;
}

.chat-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: #2a5d5d;
  font-weight: bold;
  width: 100%;
}

img {
  border-radius: 50%;
  margin-right: 15px;
  width: 40px;
  height: 40px;
  flex-shrink: 0; /* Чтобы аватарки не сжимались */
}

.chat-name {
  flex: 1; /* Растягивает название на всю оставшуюся ширину */
  text-align: left; /* Выравнивает текст влево */
  white-space: nowrap; /* Запрещает перенос текста */
  overflow: hidden; /* Скрывает лишний текст */
  text-overflow: ellipsis; /* Добавляет многоточие, если текст слишком длинный */
}


.chat-item:hover {
  background: #e9ecef;
}

/* Кнопки */
.button {
  background: #2a5d5d;
  color: white;
  border: none;
  padding: 10px 15px;
  border-radius: 5px;
  cursor: pointer;
  transition: background 0.3s;
  text-decoration: none;
  font-weight: bold;
}

.button:hover {
  background: #1e3c3c;
}

/* Позиционирование */
.search-users {
  position: absolute;
  top: 20px;
  left: 20px;
}

.profile {
  position: absolute;
  top: 20px;
  right: 20px;
}

</style>
