import React from 'react'
import { ListComponent } from './list.component'
import { getMemberCollection } from './list.repository'
import { MemberEntity } from './list.vm'

export const ListContainer: React.FC = () => {
  const [members, setMembers] = React.useState<MemberEntity[]>([])

  React.useEffect(() => {
    getMemberCollection().then((memberCollection) =>
      setMembers(memberCollection)
    )
  }, [])

  return <ListComponent members={members} />
}
