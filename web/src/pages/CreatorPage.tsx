import Creator from "../models/creator";

interface CreatorPageProps {
  creator: Creator;
}

function CreatorPage(props: CreatorPageProps) {
  return <h1>Welcome {props.creator.nickname}</h1>
}

export default CreatorPage;
