import { useEffect, useState } from 'react';
import ListProgrammers from '../../components/ListProgrammers';
import { fetchProgrammers, filterProgrammersBySkill } from '../../api';

import "./styles.css"

const Home = () => {
  const [searchSkill, setSearchSkill] = useState('');
  const [programmers, setProgrammers] = useState([]);
  const [isDeleted, setIsDeleted] = useState(false);

  useEffect(() => {
    loadData();
  }, [isDeleted]);

  const loadData = async () => {
    try {
      const res = await fetchProgrammers();
      setProgrammers(res);
      setIsDeleted(false);
    } catch (error) {
      console.error(error);
    }
  };

  async function handleSkillSearch() {
    try {
      const res = await filterProgrammersBySkill(searchSkill);
      setProgrammers(res);
    } catch(error) {
      console.error(error);
    }

    setSearchSkill('');
  }

  return (
    <>
      <h1 style={{
        textAlign: "center"
      }}>
        Search for Programmers
      </h1>

      <div style={{border: "1px solid #c9c9c937", margin: "20px"}}>
        <div className="searchContainer">
          <input type="text" placeholder="Search for a skill" value={searchSkill} onChange={(e) => {setSearchSkill(e.target.value)}} required />
          <button type="submit" onClick={handleSkillSearch} className="submitButton">Search</button>
        </div>

        <ListProgrammers programmers={programmers} setIsDeleted={setIsDeleted}/>
      </div>
    </>
  )
};

export default Home;