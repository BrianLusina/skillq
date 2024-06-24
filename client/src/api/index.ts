import axios from 'axios';

const BASE_URL= 'http://localhost:5001'

export const fetchProgrammers = async () => {
  try {
    const response = await axios.get(`${BASE_URL}/api/v1/users/`);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    throw new Error(error);
  }
};

export const deleteProgrammer = async (userId: string) => {
  try {
    await axios.delete(`${BASE_URL}/api/v1/users/delete/${userId}`);
    console.log('User deleted');
  } catch (error) {
    console.error(error);
  }
};

export const filterProgrammersBySkill = async (searchSkill: string) => {
  try {
    const response = await axios.get(`${BASE_URL}/api/v1/users/skill/${searchSkill}`);
    return response.data;
  } catch (error) {
    console.error(error);
  }
}

type imageRequest = {
  type: string;
  image: string | ArrayBuffer | null;
};

type CreateUserRequest = {
  name: string;
  email: string;
  jobTitle: string;
  image: imageRequest;
  skills: string[];
}

export const createUser = async(payload: CreateUserRequest) => {
  try {
    // Send a POST request to the back-end API to create a new user
    await axios
        .post(`${BASE_URL}/api/v1/users/`, payload)
        .then((response) => {
            console.log(response);
        })
        .catch((error) => {
            console.log(error);
        });
    
  } catch (error) {
    console.error(error)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    throw new Error(error)
  }
}