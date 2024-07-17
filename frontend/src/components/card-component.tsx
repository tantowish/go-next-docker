
interface Card {
    id: number,
    name: string,
    email: string
}

export default function CardComponent({card}: {card: Card}) {
  return (
    <div className="bg-cyan-500 p-2 text-white rounded-lg">
        <h2>{card.name}</h2>
        <h2>{card.email}</h2>
    </div>
  )
}
