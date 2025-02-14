export default function UserChatBubble({ text }: { text: string }) {
  return (
    <div className="flex flex-row-reverse">
      <div className="flex max-w-[75%] rounded-md bg-[#f3f3f3] p-4 text-sm text-gray-900">
        {text}
      </div>
    </div>
  );
}
