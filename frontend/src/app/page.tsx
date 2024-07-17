import UserInterface from "@/components/user-interface";
import Image from "next/image";

export default function Home() {
  return (
    <div className="flex flex-wrap min-h-screen justify-center items-center bg-zinc-200">
      <UserInterface />
    </div>
  );
}
