type User = {
    xid: string;
    uuid: string;
    keyId: string;
    imageUrl: string;
    name: string;
    jobTitle: string;
    email: string;
    skills: string[];
}

export type ListProgrammersProps = {
    programmers: User[];
    setIsDeleted: (e: boolean) => void;
}
