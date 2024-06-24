import { FunctionComponent } from 'react';
import { deleteProgrammer } from '../../api';
import { ListProgrammersProps } from './ListProgrammers.props';
import './listProgrammers.css';

const ListProgrammers: FunctionComponent<ListProgrammersProps> = ({programmers, setIsDeleted}) => {

  async function handleDeleteUser(userId: string) {
    deleteProgrammer(userId);
    setIsDeleted(true);

    alert("Programmer deleted successfully!");
  }

  return (
    <div>
      <ul>
        {
          (!programmers || programmers.length === 0)
          ? 
          <div className="container">
            <p>No programmers data available.</p> 
          </div>
          : 
          <>
            {
              programmers.map((user) => (
                <div key={user.uuid} className="container">
                  <div className="imageColumn">
                    <img className='userImage' src={user.imageUrl} alt="User Identity" />
                  </div>
                  <div className="infoColumn">
                    <div className="row">
                      <h2 className="userName">{user.name}</h2>
                    </div>
                    <div className="row">
                      <p className="jobTitle">
                        <span style={{fontWeight: "bold"}}>
                          Job Title:&nbsp;
                        </span>
                        {user.jobTitle}
                      </p>
                    </div>
                    <div className="row">
                      <p className="jobTitle">
                        <span style={{fontWeight: "bold"}}>
                          Email:&nbsp;
                        </span>
                        {user.email}
                      </p>
                    </div>
                    <div className="row">
                      <div className="skills">
                        {user.skills.map((skill, index) => (
                          <div key={index} className="skillBox">{skill}</div>
                        ))}
                      </div>
                    </div>
                    
                    <div className="row" style={{marginTop: "20px"}}>
                      <button type="submit" className='deleteButton' onClick={() => handleDeleteUser(user.uuid)}>Delete User</button>
                    </div>
                  </div>
                </div>
              ))
            }
          </>
        }
      </ul>
    </div>
  );
};

export default ListProgrammers;