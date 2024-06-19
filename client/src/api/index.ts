import axios from 'axios';

export const fetchProgrammers = async () => {
  try {
    const response = await axios.get('/users');
    return response.data;
  } catch (error) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    throw new Error(error);
  }
};

export const deleteProgrammer = async (userId: string) => {
  try {
    await axios.delete(`/users/delete/${userId}`);
    console.log('User deleted');
  } catch (error) {
    console.error(error);
  }
};

export const filterProgrammersBySkill = async (searchSkill: string) => {
  try {
    const response = await axios.get(`/users/skill/${searchSkill}`);
    return response.data;
  } catch (error) {
    console.error(error);
  }
}