import { ChangeEvent, FormEvent, useState } from "react";
import axios from "axios";
import "./addProgrammer.css";

const AddProgrammer = () => {
    const [name, setName] = useState<string>("");
    const [email, setEmail] = useState<string>("");
    const [jobTitle, setJobTitle] = useState<string>("");
    const [imageFile, setImageFile] = useState<File | null>(null);
    const [skills, setSkills] = useState<string[]>([]);
    const [newSkill, setNewSkill] = useState<string>("");

    function handleImageChange(e: ChangeEvent<HTMLInputElement>) {
        const {
            target: { files },
        } = e;
        if (files) {
            const file = files[0];
            setImageFile(file);
        }
    }

    function handleSkillAdd() {
        setSkills([...skills, newSkill]);
        setNewSkill("");
    }

    async function handleSubmit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault()
        if (imageFile) {
            // Convert image to a string
            const image = await readFileAsBase64(imageFile);

            // Send a POST request to the back-end API to create a new user
            axios
                .post("/users/create", { name, email, jobTitle, image, skills })
                .then((response) => {
                    console.log(response);
                })
                .catch((error) => {
                    console.log(error);
                });

            // Reset the form values
            setName("");
            setEmail("");
            setJobTitle("");
            setImageFile(null);
            setSkills([]);
            setNewSkill("");

            alert("Programmer added successfully!");
        }
    }

    // Helper function to read the image file as a base64 string
    const readFileAsBase64 = (file: File) => {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(reader.result);
            reader.onerror = (error) => reject(error);
            reader.readAsDataURL(file);
        });
    };

    return (
        <div className="formContainer">
            <form onSubmit={handleSubmit}>
                <h1 className="addProgrammerHeading">Add a new Programmer</h1>
                <div className="alignSameLine">
                    <label>Name:</label>
                    <input
                        type="text"
                        value={name}
                        onChange={(e) => {
                            setName(e.target.value);
                        }}
                        required
                    />
                </div>
                <div className="alignSameLine">
                    <label>Email:</label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => {
                            setEmail(e.target.value);
                        }}
                        required
                    />
                </div>
                <div className="alignSameLine">
                    <label>Job Title:</label>
                    <input
                        type="text"
                        value={jobTitle}
                        onChange={(e) => {
                            setJobTitle(e.target.value);
                        }}
                        required
                    />
                </div>
                <div className="alignSameLine">
                    <label>Image:</label>
                    <input
                        type="file"
                        accept="image/*"
                        onChange={handleImageChange}
                        className="imageInput"
                        required
                    />
                    <div className="tooltip">
                        Please upload a PNG image of size less than 500kb.
                    </div>
                </div>
                <>
                    <div className="alignSkillsItems">
                        <label>Skills:</label>
                        <input
                            type="text"
                            className="skillInput"
                            value={newSkill}
                            onChange={(e) => {
                                setNewSkill(e.target.value);
                            }}
                        />
                        <button
                            type="button"
                            className="skillButton"
                            onClick={handleSkillAdd}
                        >
                            Add Skill
                        </button>
                    </div>
                    {skills.length === 0 ? (
                        <div className="skillsMessage">
                            <span>Please add atleast 1 skill.</span>
                        </div>
                    ) : (
                        <div>
                            <ul>
                                {skills.map((skill, index) => (
                                    <li key={index}>{skill}</li>
                                ))}
                            </ul>
                        </div>
                    )}
                </>
                <button type="submit" className="submitButton">
                    Submit
                </button>
            </form>
        </div>
    );
};

export default AddProgrammer;
